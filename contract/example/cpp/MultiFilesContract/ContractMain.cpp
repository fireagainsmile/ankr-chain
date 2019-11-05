 #include "ContractMain.h"
 #include "InvokedClass.h"

 #include "../include/akchainlib.h"

 char* ContractMain::init() {
    print_s("TestContract::init");

    return "";
 }

void ContractMain::testFunc(const char *testStr) {
    print_s("TestContract::testFunc");
    print_s(testStr);

    char* data= "{\"Param\":\"testStr\"}";

    TrigEvent("testFunc(string)", data);
 }

int ContractMain::testFuncWithInt(const char *testStr) {
    print_s("TestContract::testFuncWithInt");
    return 1;
 }

char* ContractMain::testFuncWithString(const char *testStr) {
    print_s("TestContract::testFuncWithString");

    InvokedClass inclass;
    inclass.testFunc2();

    char* data= "[{\"index\":1,\"Name\":\"testStr\",\"ParamType\":\"string\",\"Value\":\"testStrVal\"}]";

    TrigEvent("testFunc(string)", data);

    return "";
}
