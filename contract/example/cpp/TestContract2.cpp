#include "include/contract.h"
#include "include/akchainlib.h"

class TestContract : public akchain::Contract {
public:
    [[ACTION]]char* init();
    [[ACTION]][[EVENT]] void testFunc(const char *testStr);
    [[ACTION]] int testFuncWithInt(const char *testStr);
    [[ACTION]][[EVENT]][[OWNERABLE]] char *testFuncWithString(const char *testStr);
 };

 char* TestContract::init() {
    print_s("TestContract::init");

    return "";
 }

void TestContract::testFunc(const char *testStr) {
    print_s("TestContract::testFunc");
    print_s(testStr);

    char* data= "{\"Param\":\"testStr\"}";

    TrigEvent("testFunc(string)", data);
 }

int TestContract::testFuncWithInt(const char *testStr) {
    print_s("TestContract::testFuncWithInt");
    return 1;
 }

char* TestContract::testFuncWithString(const char *testStr) {
    print_s("TestContract::testFuncWithString");

    char* data= "[{\"index\":1,\"Name\":\"testStr\",\"ParamType\":\"string\",\"Value\":\"testStrVal\"}]";

    TrigEvent("testFunc(string)", data);

    return "";
}

//auto generated

extern "C" {

    EXPORT char *init() {
        TestContract tc;
        return tc.init();
    }

    EXPORT void testFunc(const char *testStr) {
         TestContract tc;
         tc.testFunc(testStr);
    }

    EXPORT int testFuncWithInt(const char *testStr) {
         TestContract tc;
         return tc.testFuncWithInt(testStr);
    }

    EXPORT char *testFuncWithString(const char *testStr) {
         TestContract tc;
         return tc.testFuncWithString(testStr);
    }
}








