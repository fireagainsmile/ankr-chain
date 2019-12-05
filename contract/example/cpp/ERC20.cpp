#include "include/contract.h"
#include "include/akchainlib.h"

class ERC20 : public akchain::Contract {
public:
    [[ACTION]]char* init();
    [[ACTION]]char* Name();
    [[ACTION]]char* Symbol();
    [[ACTION]]int Decimals();
    [[ACTION]]char* TotalSupply();
    [[ACTION]]char* BalanceOf(const char* addr);
    [[ACTION]][[EVENT]]int Transfer(const char* toAddr, const char* amount);
    [[ACTION]][[EVENT]]int TransferFrom(const char* fromAddr, const char* toAddr, const char* amount);
    [[ACTION]][[EVENT]]int Approve(const char* spenderAddr, const char* amount);
    [[ACTION]]char* AllowanceERC20(const char* ownerAddr, const char* spenderAddr);
    [[ACTION]][[EVENT]]int IncreaseApproval(const char* spenderAddr, const char* addedAmount);
    [[ACTION]][[EVENT]]int DecreaseApproval(const char* spenderAddr, const char* addedAmount);
 };

char* ERC20::init() {
    char* cAddr = ContractAddr();
    char* senderAddr = SenderAddr();
    CreateCurrency(Symbol(), Decimals(), TotalSupply());
    BuildCurrencyCAddrMap(Symbol(), cAddr);
    SetBalance(senderAddr, Symbol(),TotalSupply());

    return "";
}

char* ERC20::Name() {
    return "ERC20";
}
char* ERC20::Symbol() {
    return "TESTCOIN";
}

int ERC20::Decimals() {
    return 18;
}
char* ERC20::TotalSupply() {
    return "1000000000000000000000000000";
}

char* ERC20::BalanceOf(const char* addr) {
    return Balance(addr, "TESTCOIN");
}

int ERC20::Transfer(const char* toAddr, const char* amount) {
    if (strcmp(toAddr, "") == 0 || strcmp(amount, "") == 0){
        return -1;
    }
    char* senderAddr = SenderAddr();
    char* balSender = BalanceOf(senderAddr);
    char* balTo     = BalanceOf(toAddr);

     if (balSender == nullptr || BigIntCmp(balSender,amount) <= 0) {
        return -1;
     }

     if (balTo == nullptr || strcmp(balTo, "") == 0) {
        balTo = "0";
     }

    balSender = BigIntSub(balSender, amount);
    balTo     = BigIntAdd(balTo, amount);

    SetBalance(senderAddr, "TESTCOIN", balSender);
    SetBalance(toAddr, "TESTCOIN", balTo);

    char* jsonArg = "[{\"index\":1,\"Name\":\"toAddr\",\"ParamType\":\"string\",\"Value\":\"toAddrVal\"},"
    		         "{\"index\":2,\"Name\":\"amount\",\"ParamType\":\"string\",\"Value\":\"amontVal\"}]";

    TrigEvent("transfer(string, string)", jsonArg);

    return 0;
}

int ERC20::TransferFrom(const char* fromAddr, const char* toAddr, const char* amount) {
     if (strcmp(fromAddr, "") == 0 || strcmp(toAddr, "") == 0 || strcmp(amount, "") == 0){
        return -1;
     }

    char* balFrom = BalanceOf(fromAddr);
    char* balTo   = BalanceOf(toAddr);

    if (balFrom == nullptr || BigIntCmp(balFrom,amount) <= 0) {
        return -1;
    }

    if (balTo == nullptr || strcmp(balTo, "") == 0) {
        balTo = "0";
    }

    balFrom = BigIntSub(balFrom, amount);
    balTo   = BigIntAdd(balTo, amount);

    SetBalance(fromAddr, "TESTCOIN", balFrom);
    SetBalance(toAddr, "TESTCOIN", balTo);

     char* jsonArg = "[{\"index\":1,\"Name\":\"fromAddr\",\"ParamType\":\"string\",\"Value\":\"fromAddrVal\"},"
                      "{\"index\":1,\"Name\":\"toAddr\",\"ParamType\":\"string\",\"Value\":\"toAddrVal\"},"
        		      "{\"index\":2,\"Name\":\"amount\",\"ParamType\":\"string\",\"Value\":\"amontVal\"}]";

     TrigEvent("transferFrom(string, string, string)", jsonArg);

    return 0;
}

int ERC20::Approve(const char* spenderAddr, const char* amount) {
    char* senderAddr = SenderAddr();
    int iRtn = SetAllowance(senderAddr, spenderAddr, "TESTCOIN", amount);

    char* jsonArg = "[{\"index\":1,\"Name\":\"spenderAddr\",\"ParamType\":\"string\",\"Value\":\"spenderAddrVal\"},"
                     "{\"index\":2,\"Name\":\"amount\",\"ParamType\":\"string\",\"Value\":\"amontVal\"}]";

    TrigEvent("approve(string, string)", jsonArg);

    return iRtn;
}

char* ERC20::AllowanceERC20(const char* ownerAddr, const char* spenderAddr) {
    return Allowance(ownerAddr, spenderAddr, "TESTCOIN");
}

int ERC20::IncreaseApproval(const char* spenderAddr, const char* addedAmount) {
    char* senderAddr = SenderAddr();
    char* curAllow = Allowance(senderAddr, spenderAddr, "TESTCOIN");
    char* allow = BigIntAdd(curAllow, addedAmount);
    int iRtn = SetAllowance(senderAddr, spenderAddr, "TESTCOIN", allow);

    char* jsonArg = "[{\"index\":1,\"Name\":\"spenderAddr\",\"ParamType\":\"string\",\"Value\":\"spenderAddrVal\"},"
                     "{\"index\":2,\"Name\":\"addedAmount\",\"ParamType\":\"string\",\"Value\":\"addedAmountVal\"}]";

    TrigEvent("increaseApproval(string, string)", jsonArg);

    return iRtn;
}

int ERC20::DecreaseApproval(const char* spenderAddr, const char* subtractedAmount) {
    char* senderAddr = SenderAddr();
    char* curAllow = Allowance(senderAddr, spenderAddr, "TESTCOIN");
    char* allow = BigIntSub(curAllow, subtractedAmount);
    int iRtn = SetAllowance(senderAddr, spenderAddr, "TESTCOIN", allow);

    char* jsonArg = "[{\"index\":1,\"Name\":\"spenderAddr\",\"ParamType\":\"string\",\"Value\":\"spenderAddrVal\"},"
                     "{\"index\":2,\"Name\":\"subtractedAmount\",\"ParamType\":\"string\",\"Value\":\"subtractedAmountVal\"}]";

    TrigEvent("decreaseApproval(string, string)", jsonArg);

    return iRtn;
}
