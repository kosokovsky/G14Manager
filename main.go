package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zllovesuki/ROGManager/controller"
	"github.com/zllovesuki/ROGManager/system/battery"
	"github.com/zllovesuki/ROGManager/system/persist"
	"github.com/zllovesuki/ROGManager/system/thermal"

	"github.com/google/shlex"
)

var defaultCommandWithArgs = "Taskmgr.exe"

func main() {

	rogRemap := flag.String("rog", defaultCommandWithArgs, "which program to launch when the ROG Key is pressed")

	flag.Parse()

	commandWithArgs, err := shlex.Split(*rogRemap)
	if err != nil {
		log.Fatalln(err)
	}

	config, _ := persist.NewRegistryHelper()

	powercfg, err := thermal.NewPowerCfg()
	if err != nil {
		log.Fatalln(err)
	}
	config.Register(powercfg)

	// TODO: allow user to specify profiles
	thermalCfg := thermal.Config{
		PowerCfg: powercfg,
		Profiles: thermal.GetDefaultThermalProfiles(),
	}

	profile, err := thermal.NewThermal(thermalCfg)
	if err != nil {
		log.Fatalln(err)
	}
	config.Register(profile)

	// TODO: allow user to change the charge limit
	battery, err := battery.NewChargeLimit()
	if err != nil {
		log.Fatalln(err)
	}
	config.Register(battery)

	// load configs from registry and try to reapply
	if err := config.Load(); err != nil {
		log.Fatalln(err)
	}
	/*if err := config.Apply(); err != nil {
		log.Fatalln(err)
	}*/

	control, err := controller.NewController(controller.Config{
		Thermal:  profile,
		Registry: config,
		ROGKey:   commandWithArgs,
	})

	if err != nil {
		log.Fatalln(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sigc
		log.Printf("graceful shutdown on %+v\n", s)
		control.Shutdown()
	}()

	os.Exit(control.Run())
}
