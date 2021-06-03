//
// Created by xmmmmmovo on 2021/6/3.
//

#ifndef EMERALDDB_APPLE_HPP
#define EMERALDDB_APPLE_HPP

#ifdef __APPLE__
int gettid() {
    uint64_t tid;
    pthread_threadid_np(NULL, &tid);
    return tid;
}
#endif

#endif // EMERALDDB_APPLE_HPP
