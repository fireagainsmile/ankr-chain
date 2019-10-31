#ifndef CONTRACT_H_
#define CONTRACT_H_

namespace akchain {

class Contract {
public:
    Contract(){}
    virtual ~Contract(){}
    virtual char* init() = 0;
    virtual char* actionEntry(const char* action_name, const char *args) { return ""; }
};

}

#endif /*CONTRACT_H_*/