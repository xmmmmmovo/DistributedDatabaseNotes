//
// Created by xmmmmmovo on 2021/6/2.
//

#include "ossPrimitiveFileOp.hpp"
#include "core.hpp"

ossPrimitiveFileOp::ossPrimitiveFileOp() {
    _fileHandle = OSS_INVALID_HANDLE_FD_VALUE;
    _bIsStdout  = false;
}
bool ossPrimitiveFileOp::isValid() {
    return (OSS_INVALID_HANDLE_FD_VALUE != _fileHandle);
}
void ossPrimitiveFileOp::Close() {
    if (isValid() && (!_bIsStdout)) {
        oss_close(_fileHandle);
        _fileHandle = OSS_INVALID_HANDLE_FD_VALUE;
    }
}
