#pragma once
enum class MatrixStatus
{
	NO_CONTROLLER				= 99,
	MATRIX_OPERATION_SUCCESSFUL = 0,
	MATRIX_NO_HANDLE			= 1,
	MATRIX_VECTOR_INPUT_ERROR	= 2,
	MATRIX_DRAW_ERROR			= 3
};