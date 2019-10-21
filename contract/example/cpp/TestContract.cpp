#include "include/contract.h"
#include "include/akchainlib.h"

class TestContract : public akchain::Contract {
public:
    char* init();
    void testFunc(const char *testStr);
    int testFuncWithInt(const char *testStr);
    char *testFuncWithString(const char *testStr);
    char* actionEntry(const char* action_name, const char *args);
 };

char* TestContract::init() {
    return "";
}

void TestContract::testFunc(const char *testStr) {
     print_s(testStr);
     int i = strcmp("sdf", testStr);
     print_i(i);
}

int TestContract::testFuncWithInt(const char *testStr) {
    print_s(testStr);
    return 100;
}

char *TestContract::testFuncWithString(const char *testStr) {
     print_s(testStr);
     return "testFuncWithString sucess";
}

char* TestContract::actionEntry(const char* action, const char *args) {
    int index = JsonObjectIndex(args);

    INVOKE_ACTION("init",
        TestContract::init();
    );

    INVOKE_ACTION("testFunc",
        const char *argStr = JsonGetString(index, "testStr");
        testFunc(argStr);
    );

    INVOKE_ACTION("testFuncWithInt",
        const char *argStr = JsonGetString(index, "testStr");
        testFuncWithInt(argStr);
    );

   INVOKE_ACTION("testFuncWithString",
        const char *argStr = JsonGetString(index, "testStr");
        char *rtnStr = testFuncWithString(argStr);
        int index = JsonCreateObject();
        JsonPutString(index, "testFuncWithString", rtnStr);
        return JsonToString(index);
    );

   return "";
}

extern "C" {

EXPORT char *ContractEntry(char *action, char *args) {
        TestContract tc;
         print_s(action);
         print_s(args);
        return tc.actionEntry(action, args);
       }

}













