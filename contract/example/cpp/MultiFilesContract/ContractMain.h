#include "../include/contract.h"

#ifndef CONTRACTMAIN_H_
#define CONTRACTMAIN_H_

class ContractMain : public akchain::Contract {
public:
    [[ACTION]]char* init();
    [[ACTION]][[EVENT]] void testFunc(const char *testStr);
    [[ACTION]] int testFuncWithInt(const char *testStr);
    [[ACTION]][[EVENT]]char *testFuncWithString(const char *testStr);
 };

#endif/*CONTRACTMAIN_H_*/
