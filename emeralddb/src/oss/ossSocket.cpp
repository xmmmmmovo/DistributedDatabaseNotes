//
// Created by xmmmmmovo on 2021/6/2.
//

#include "ossSocket.hpp"

_ossSocket::_ossSocket() {
    _init    = false;
    _fd      = 0;
    _timeout = 0;
    memset(&_sockAddress, 0, sizeof(sockaddr_in));
    memset(&_peerAddress, 0, sizeof(sockaddr_in));
    _peerAddressLen = sizeof(_peerAddress);
    _addressLen     = sizeof(_sockAddress);
}

_ossSocket::_ossSocket(unsigned int port, int timeout) {
    _init = false;
    _fd   = 0;
    memset(&_sockAddress, 0, sizeof(sockaddr_in));
    memset(&_peerAddress, 0, sizeof(sockaddr_in));
    _peerAddressLen              = sizeof(_peerAddress);
    _peerAddress.sin_family      = AF_INET;
    _sockAddress.sin_addr.s_addr = htonl(INADDR_ANY);
    _sockAddress.sin_port        = htons(port);
    _addressLen                  = sizeof(_sockAddress);
}
_ossSocket::_ossSocket(const char *pHostname, unsigned int port,
                       int timeout) {
    hostent *hp;
    _init    = false;
    _timeout = timeout;
    _fd      = 0;
    memset(&_sockAddress, 0, sizeof(sockaddr_in));
    memset(&_peerAddress, 0, sizeof(sockaddr_in));
    _peerAddressLen         = sizeof(_peerAddress);
    _sockAddress.sin_family = AF_INET;
    if ((hp = gethostbyname(pHostname))) {
        _sockAddress.sin_addr.s_addr = *((int *)hp->h_addr_list[0]);
    } else {
        _sockAddress.sin_addr.s_addr = inet_addr(pHostname);
    }
    _sockAddress.sin_port = htons(port);
    _addressLen           = sizeof(_sockAddress);
}

// Create from a existing socket
_ossSocket::_ossSocket(int *sock, int timeout) {
    int rc      = EDB_OK;
    _fd         = *sock;
    _init       = true;
    _timeout    = timeout;
    _addressLen = sizeof(_sockAddress);
    memset(&_peerAddress, 0, sizeof(sockaddr_in));
    _peerAddressLen = sizeof(_peerAddress);
    rc = getsockname(_fd, (sockaddr *)&_sockAddress, &_addressLen);
    if (rc) {
        printf("Failed to get sock name, error = %d",
               SOCKET_GETLASTERROR);
        _init = false;
    } else {
        rc = getpeername(_fd, (sockaddr *)&_peerAddress,
                         &_peerAddressLen);
        printf("Failed to get peer name, error = %d",
               SOCKET_GETLASTERROR);
    }
}
int _ossSocket::initSocket() {
    int rc = EDB_OK;
    if (_init) {
        goto done;
    }
    memset(&_peerAddress, 0, sizeof(sockaddr_in));
    _peerAddressLen = sizeof(_peerAddress);
    if (-1 == _fd) {
        printf("Failed to init socket error = %d",
               SOCKET_GETLASTERROR);
        rc = EDB_NETWORK;
        goto error;
    }
    _init = true;
    setTimeout(_timeout);
done:
    return rc;
error:
    goto done;
}
int _ossSocket::setSocketLi(int lOnOff, int linger) {
    int           rc      = EDB_OK;
    struct linger _linger = {};
    _linger.l_onoff       = lOnOff;
    _linger.l_linger      = linger;
    rc                    = setsockopt(_fd, SOL_SOCKET, SO_LINGER,
                    (const char *)&_linger, sizeof(_linger));
    return rc;
}
void _ossSocket::setAddress(const char * pHostname,
                            unsigned int port) {
    hostent *hp;
    memset(&_sockAddress, 0, sizeof(sockaddr_in));
    memset(&_peerAddress, 0, sizeof(sockaddr_in));
    _peerAddressLen = sizeof(_peerAddress);
    if ((hp = gethostbyname(pHostname))) {
        _sockAddress.sin_addr.s_addr = *((int *)hp->h_addr_list[0]);
    } else {
        _sockAddress.sin_addr.s_addr = inet_addr(pHostname);
    }

    _sockAddress.sin_port = htons(port);
    _addressLen           = sizeof(_sockAddress);
}
int _ossSocket::bind_listen() {
    int rc   = EDB_OK;
    int temp = 1;
    rc = setsockopt(_fd, SOL_SOCKET, SO_REUSEADDR, (char *)&temp,
                    sizeof(int));
    if (rc) {
        printf("error=%d", SOCKET_GETLASTERROR);
    }
    rc = ::bind(_fd, (sockaddr *)&_sockAddress, _addressLen);
    if (rc) {
        printf("Failed to bind socket, rc = %d", SOCKET_GETLASTERROR);
        rc = EDB_NETWORK;
        goto error;
    }
    rc = listen(_fd, SOMAXCONN);
    if (rc) {
        printf("Failed to listen socket, rc = %d",
               SOCKET_GETLASTERROR);
        rc = EDB_NETWORK;
        goto error;
    }

done:
    return rc;
error:
    close();
    goto done;
}
int _ossSocket::send(const char *pMsg, int len, int timeout,
                     int flags) {
    int     rc    = EDB_OK;
    int     maxFD = _fd;
    timeval maxSelectTime;
    fd_set  fds;

    maxSelectTime.tv_sec  = timeout / 1000000;
    maxSelectTime.tv_usec = timeout % 1000000;
    if (0 == len) {
        return EDB_OK;
    }

    while (true) {
        FD_ZERO(&fds);
        FD_SET(_fd, &fds);
        rc = select(maxFD + 1, NULL, &fds, NULL,
                    timeout >= 0 ? &maxSelectTime : NULL);

        if (0 == rc) {
            rc = EDB_TIMEOUT;
            goto done;
        }

        if (0 > rc)
        {
            rc = SOCKET_GETLASTERROR;
            if (EINTR == rc)
            {
                continue;
            }
            printf("Failed to select from socket, rc = %d", rc);
            rc = EDB_NETWORK;
            goto error;
        }

        if (FD_ISSET(_fd, &fds))
        {
           break;
        }
        
        
        
    }
}
