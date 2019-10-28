#include "include/contract.h"
#include "include/akchainlib.h"

class ERC20 : public akchain::Contract {
public:
    [[ACTION]]char* init();
    [[ACTION]]char* name();
    [[ACTION]]char* symbol();
    [[ACTION]]int decimals();
    [[ACTION]]char* totalSupply();
    [[ACTION]]char* balanceOf(const char* addr);
    [[ACTION]][[EVENT]]int transfer(const char* toAddr, const char* amount);
    [[ACTION]][[EVENT]]int transferFrom(const char* fromAddr, const char* toAddr, const char* amount);
    [[ACTION]][[EVENT]]int approve(const char* spenderAddr, const char* amount);
    [[ACTION]]char* allowance(const char* ownerAddr, const char* spenderAddr);
    [[ACTION]][[EVENT]]int increaseApproval(const char* spenderAddr, const char* addedAmount);
    [[ACTION]][[EVENT]]int decreaseApproval(const char* spenderAddr, const char* addedAmount);
 };

char* ERC20::init() {
    CreateCurrency("TESTCOIN", 18);
    char* cAddr = ContractAddr();
    BuildCurrencyCAddrMap("TESTCOIN", cAddr);
}

char* ERC20::name() {
    return "ERC20";
}
char* ERC20::symbol() {
    return "TESTCOIN";
}

int ERC20::decimals() {
    return 18;
}
char* ERC20::totalSupply() {
    return "1000000000000000000000000000";
}

char* ERC20::balanceOf(const char* addr) {
    return Balance(addr, "TESTCOIN");
}

int ERC20::transfer(const char* toAddr, const char* amount) {
    if (strcmp(toAddr, "") == 0 || strcmp(amount, "") == 0){
        return -1;
    }
    char* senderAddr = SenderAddr();
    char* balSender = balanceOf(senderAddr);
    char* balTo     = balanceOf(toAddr);

    if (balSender == NULL || BigCmp(balSender,amount) <= 0) {
        return -1;
    }

    balSender = BigSub(balSender, amount);
    balTo     = BigAdd(balTo, amount);

    SetBalance(senderAddr, "TESTCOIN", balSender);
    SetBalance(toAddr, "TESTCOIN", balTo);

    return 0;
}

int ERC20::transferFrom(const char* fromAddr, const char* toAddr, const char* amount) {
     if (strcmp(fromAddr, "") == 0 || strcmp(toAddr, "") == 0 || strcmp(amount, "") == 0){
        return -1;
     }

    char* balFrom = balanceOf(fromAddr);
    char* balTo   = balanceOf(toAddr);

    if (balFrom == NULL || BigCmp(balFrom,amount) <= 0) {
        return -1;
    }

    balFrom = BigSub(balFrom, amount);
    balTo   = BigAdd(balTo, amount);

    SetBalance(fromAddr, "TESTCOIN", balFrom);
    SetBalance(toAddr, "TESTCOIN", balTo);
}

int ERC20::approve(const char* spenderAddr, const char* amount) {
    char* senderAddr = SenderAddr();
    return SetAllowance(senderAddr, spenderAddr, "TESTCOIN", amount)
}

char* ERC20::allowance(const char* ownerAddr, const char* spenderAddr) {
    return Allowance(ownerAddr, spenderAddr, "TESTCOIN")
}

int ERC20::increaseApproval(const char* spenderAddr, const char* addedAmount) {
    char* senderAddr = SenderAddr();
    char* curAllow = Allowance(senderAddr, spenderAddr, "TESTCOIN");
    char* allow = BigAdd(curAllow, addedAmount);
    return SetAllowance(senderAddr, spenderAddr, "TESTCOIN", allow)
}

int ERC20::decreaseApproval(const char* spenderAddr, const char* subtractedAmount) {
    char* senderAddr = SenderAddr();
    char* curAllow = Allowance(senderAddr, spenderAddr, "TESTCOIN");
    char* allow = BigSub(curAllow, subtractedAmount);
    return SetAllowance(senderAddr, spenderAddr, "TESTCOIN", allow)
}
