#include <stdio.h>
#include <iostream>
#include <assert.h>
#include "circom.hpp"
#include "calcwit.hpp"
void NonEqual_0_create(uint soffset,uint coffset,Circom_CalcWit* ctx,std::string componentName,uint componentFather);
void NonEqual_0_run(uint ctx_index,Circom_CalcWit* ctx);
void Distinct_1_create(uint soffset,uint coffset,Circom_CalcWit* ctx,std::string componentName,uint componentFather);
void Distinct_1_run(uint ctx_index,Circom_CalcWit* ctx);
void Bits4_2_create(uint soffset,uint coffset,Circom_CalcWit* ctx,std::string componentName,uint componentFather);
void Bits4_2_run(uint ctx_index,Circom_CalcWit* ctx);
void OneToNine_3_create(uint soffset,uint coffset,Circom_CalcWit* ctx,std::string componentName,uint componentFather);
void OneToNine_3_run(uint ctx_index,Circom_CalcWit* ctx);
void Sudoku_4_create(uint soffset,uint coffset,Circom_CalcWit* ctx,std::string componentName,uint componentFather);
void Sudoku_4_run(uint ctx_index,Circom_CalcWit* ctx);
Circom_TemplateFunction _functionTable[5] = { 
NonEqual_0_run,
Distinct_1_run,
Bits4_2_run,
OneToNine_3_run,
Sudoku_4_run };
Circom_TemplateFunction _functionTableParallel[5] = { 
NULL,
NULL,
NULL,
NULL,
NULL };
uint get_main_input_signal_start() {return 1;}

uint get_main_input_signal_no() {return 162;}

uint get_total_signal_no() {return 2107;}

uint get_number_of_components() {return 577;}

uint get_size_of_input_hashmap() {return 256;}

uint get_size_of_witness() {return 973;}

uint get_size_of_constants() {return 6;}

uint get_size_of_io_map() {return 0;}

void release_memory_component(Circom_CalcWit* ctx, uint pos) {{

if (pos != 0){{

if(ctx->componentMemory[pos].subcomponents)
delete []ctx->componentMemory[pos].subcomponents;

if(ctx->componentMemory[pos].subcomponentsParallel)
delete []ctx->componentMemory[pos].subcomponentsParallel;

if(ctx->componentMemory[pos].outputIsSet)
delete []ctx->componentMemory[pos].outputIsSet;

if(ctx->componentMemory[pos].mutexes)
delete []ctx->componentMemory[pos].mutexes;

if(ctx->componentMemory[pos].cvs)
delete []ctx->componentMemory[pos].cvs;

if(ctx->componentMemory[pos].sbct)
delete []ctx->componentMemory[pos].sbct;

}}


}}


// function declarations
// template declarations
void NonEqual_0_create(uint soffset,uint coffset,Circom_CalcWit* ctx,std::string componentName,uint componentFather){
ctx->componentMemory[coffset].templateId = 0;
ctx->componentMemory[coffset].templateName = "NonEqual";
ctx->componentMemory[coffset].signalStart = soffset;
ctx->componentMemory[coffset].inputCounter = 2;
ctx->componentMemory[coffset].componentName = componentName;
ctx->componentMemory[coffset].idFather = componentFather;
ctx->componentMemory[coffset].subcomponents = new uint[0];
}

void NonEqual_0_run(uint ctx_index,Circom_CalcWit* ctx){
FrElement* signalValues = ctx->signalValues;
u64 mySignalStart = ctx->componentMemory[ctx_index].signalStart;
std::string myTemplateName = ctx->componentMemory[ctx_index].templateName;
std::string myComponentName = ctx->componentMemory[ctx_index].componentName;
u64 myFather = ctx->componentMemory[ctx_index].idFather;
u64 myId = ctx_index;
u32* mySubcomponents = ctx->componentMemory[ctx_index].subcomponents;
bool* mySubcomponentsParallel = ctx->componentMemory[ctx_index].subcomponentsParallel;
FrElement* circuitConstants = ctx->circuitConstants;
std::string* listOfTemplateMessages = ctx->listOfTemplateMessages;
FrElement expaux[6];
FrElement lvar[0];
uint sub_component_aux;
uint index_multiple_eq;
{
PFrElement aux_dest = &signalValues[mySignalStart + 2];
// load src
Fr_sub(&expaux[2],&signalValues[mySignalStart + 0],&signalValues[mySignalStart + 1]); // line circom 9
Fr_div(&expaux[0],&circuitConstants[0],&expaux[2]); // line circom 9
// end load src
Fr_copy(aux_dest,&expaux[0]);
}
Fr_sub(&expaux[3],&signalValues[mySignalStart + 0],&signalValues[mySignalStart + 1]); // line circom 10
Fr_mul(&expaux[1],&signalValues[mySignalStart + 2],&expaux[3]); // line circom 10
Fr_eq(&expaux[0],&expaux[1],&circuitConstants[0]); // line circom 10
if (!Fr_isTrue(&expaux[0])) std::cout << "Failed assert in template/function " << myTemplateName << " line 10. " <<  "Followed trace of components: " << ctx->getTrace(myId) << std::endl;
assert(Fr_isTrue(&expaux[0]));
for (uint i = 0; i < 0; i++){
uint index_subc = ctx->componentMemory[ctx_index].subcomponents[i];
if (index_subc != 0)release_memory_component(ctx,index_subc);
}
}

void Distinct_1_create(uint soffset,uint coffset,Circom_CalcWit* ctx,std::string componentName,uint componentFather){
ctx->componentMemory[coffset].templateId = 1;
ctx->componentMemory[coffset].templateName = "Distinct";
ctx->componentMemory[coffset].signalStart = soffset;
ctx->componentMemory[coffset].inputCounter = 9;
ctx->componentMemory[coffset].componentName = componentName;
ctx->componentMemory[coffset].idFather = componentFather;
ctx->componentMemory[coffset].subcomponents = new uint[81]{0};
}

void Distinct_1_run(uint ctx_index,Circom_CalcWit* ctx){
FrElement* signalValues = ctx->signalValues;
u64 mySignalStart = ctx->componentMemory[ctx_index].signalStart;
std::string myTemplateName = ctx->componentMemory[ctx_index].templateName;
std::string myComponentName = ctx->componentMemory[ctx_index].componentName;
u64 myFather = ctx->componentMemory[ctx_index].idFather;
u64 myId = ctx_index;
u32* mySubcomponents = ctx->componentMemory[ctx_index].subcomponents;
bool* mySubcomponentsParallel = ctx->componentMemory[ctx_index].subcomponentsParallel;
FrElement* circuitConstants = ctx->circuitConstants;
std::string* listOfTemplateMessages = ctx->listOfTemplateMessages;
FrElement expaux[3];
FrElement lvar[3];
uint sub_component_aux;
uint index_multiple_eq;
{
PFrElement aux_dest = &lvar[0];
// load src
// end load src
Fr_copy(aux_dest,&circuitConstants[1]);
}
{
uint aux_create = 0;
int aux_cmp_num = 0+ctx_index+1;
uint csoffset = mySignalStart+9;
uint aux_dimensions[2] = {9,9};
uint aux_positions [36]= {9,18,19,27,28,29,36,37,38,39,45,46,47,48,49,54,55,56,57,58,59,63,64,65,66,67,68,69,72,73,74,75,76,77,78,79};
for (uint i_aux = 0; i_aux < 36; i_aux++) {
uint i = aux_positions[i_aux];
std::string new_cmp_name = "nonEqual"+ctx->generate_position_array(aux_dimensions, 2, i);
NonEqual_0_create(csoffset,aux_cmp_num,ctx,new_cmp_name,myId);
mySubcomponents[aux_create+i] = aux_cmp_num;
csoffset += 3 ;
aux_cmp_num += 1;
}
}
{
PFrElement aux_dest = &lvar[1];
// load src
// end load src
Fr_copy(aux_dest,&circuitConstants[2]);
}
Fr_lt(&expaux[0],&lvar[1],&circuitConstants[1]); // line circom 16
while(Fr_isTrue(&expaux[0])){
{
PFrElement aux_dest = &lvar[2];
// load src
// end load src
Fr_copy(aux_dest,&circuitConstants[2]);
}
Fr_lt(&expaux[0],&lvar[2],&lvar[1]); // line circom 17
while(Fr_isTrue(&expaux[0])){
{
uint cmp_index_ref = (((9 * Fr_toInt(&lvar[1])) + (1 * Fr_toInt(&lvar[2]))) + 0);
{
PFrElement aux_dest = &ctx->signalValues[ctx->componentMemory[mySubcomponents[cmp_index_ref]].signalStart + 0];
// load src
// end load src
Fr_copy(aux_dest,&signalValues[mySignalStart + ((1 * Fr_toInt(&lvar[1])) + 0)]);
}
// run sub component if needed
if(!(ctx->componentMemory[mySubcomponents[cmp_index_ref]].inputCounter -= 1)){
NonEqual_0_run(mySubcomponents[cmp_index_ref],ctx);

}
}
{
uint cmp_index_ref = (((9 * Fr_toInt(&lvar[1])) + (1 * Fr_toInt(&lvar[2]))) + 0);
{
PFrElement aux_dest = &ctx->signalValues[ctx->componentMemory[mySubcomponents[cmp_index_ref]].signalStart + 1];
// load src
// end load src
Fr_copy(aux_dest,&signalValues[mySignalStart + ((1 * Fr_toInt(&lvar[2])) + 0)]);
}
// run sub component if needed
if(!(ctx->componentMemory[mySubcomponents[cmp_index_ref]].inputCounter -= 1)){
NonEqual_0_run(mySubcomponents[cmp_index_ref],ctx);

}
}
{
PFrElement aux_dest = &lvar[2];
// load src
Fr_add(&expaux[0],&lvar[2],&circuitConstants[0]); // line circom 17
// end load src
Fr_copy(aux_dest,&expaux[0]);
}
Fr_lt(&expaux[0],&lvar[2],&lvar[1]); // line circom 17
}
{
PFrElement aux_dest = &lvar[1];
// load src
Fr_add(&expaux[0],&lvar[1],&circuitConstants[0]); // line circom 16
// end load src
Fr_copy(aux_dest,&expaux[0]);
}
Fr_lt(&expaux[0],&lvar[1],&circuitConstants[1]); // line circom 16
}
for (uint i = 0; i < 81; i++){
uint index_subc = ctx->componentMemory[ctx_index].subcomponents[i];
if (index_subc != 0)release_memory_component(ctx,index_subc);
}
}

void Bits4_2_create(uint soffset,uint coffset,Circom_CalcWit* ctx,std::string componentName,uint componentFather){
ctx->componentMemory[coffset].templateId = 2;
ctx->componentMemory[coffset].templateName = "Bits4";
ctx->componentMemory[coffset].signalStart = soffset;
ctx->componentMemory[coffset].inputCounter = 1;
ctx->componentMemory[coffset].componentName = componentName;
ctx->componentMemory[coffset].idFather = componentFather;
ctx->componentMemory[coffset].subcomponents = new uint[0];
}

void Bits4_2_run(uint ctx_index,Circom_CalcWit* ctx){
FrElement* signalValues = ctx->signalValues;
u64 mySignalStart = ctx->componentMemory[ctx_index].signalStart;
std::string myTemplateName = ctx->componentMemory[ctx_index].templateName;
std::string myComponentName = ctx->componentMemory[ctx_index].componentName;
u64 myFather = ctx->componentMemory[ctx_index].idFather;
u64 myId = ctx_index;
u32* mySubcomponents = ctx->componentMemory[ctx_index].subcomponents;
bool* mySubcomponentsParallel = ctx->componentMemory[ctx_index].subcomponentsParallel;
FrElement* circuitConstants = ctx->circuitConstants;
std::string* listOfTemplateMessages = ctx->listOfTemplateMessages;
FrElement expaux[6];
FrElement lvar[2];
uint sub_component_aux;
uint index_multiple_eq;
{
PFrElement aux_dest = &lvar[0];
// load src
// end load src
Fr_copy(aux_dest,&circuitConstants[2]);
}
{
PFrElement aux_dest = &lvar[1];
// load src
// end load src
Fr_copy(aux_dest,&circuitConstants[2]);
}
Fr_lt(&expaux[0],&lvar[1],&circuitConstants[3]); // line circom 30
while(Fr_isTrue(&expaux[0])){
{
PFrElement aux_dest = &signalValues[mySignalStart + ((1 * Fr_toInt(&lvar[1])) + 1)];
// load src
Fr_shr(&expaux[1],&signalValues[mySignalStart + 0],&lvar[1]); // line circom 31
Fr_band(&expaux[0],&expaux[1],&circuitConstants[0]); // line circom 31
// end load src
Fr_copy(aux_dest,&expaux[0]);
}
Fr_sub(&expaux[3],&signalValues[mySignalStart + ((1 * Fr_toInt(&lvar[1])) + 1)],&circuitConstants[0]); // line circom 32
Fr_mul(&expaux[1],&signalValues[mySignalStart + ((1 * Fr_toInt(&lvar[1])) + 1)],&expaux[3]); // line circom 32
Fr_eq(&expaux[0],&expaux[1],&circuitConstants[2]); // line circom 32
if (!Fr_isTrue(&expaux[0])) std::cout << "Failed assert in template/function " << myTemplateName << " line 32. " <<  "Followed trace of components: " << ctx->getTrace(myId) << std::endl;
assert(Fr_isTrue(&expaux[0]));
{
PFrElement aux_dest = &lvar[0];
// load src
Fr_pow(&expaux[3],&circuitConstants[4],&lvar[1]); // line circom 33
Fr_mul(&expaux[2],&expaux[3],&signalValues[mySignalStart + ((1 * Fr_toInt(&lvar[1])) + 1)]); // line circom 33
Fr_add(&expaux[0],&lvar[0],&expaux[2]); // line circom 33
// end load src
Fr_copy(aux_dest,&expaux[0]);
}
{
PFrElement aux_dest = &lvar[1];
// load src
Fr_add(&expaux[0],&lvar[1],&circuitConstants[0]); // line circom 30
// end load src
Fr_copy(aux_dest,&expaux[0]);
}
Fr_lt(&expaux[0],&lvar[1],&circuitConstants[3]); // line circom 30
}
Fr_eq(&expaux[0],&lvar[0],&signalValues[mySignalStart + 0]); // line circom 35
if (!Fr_isTrue(&expaux[0])) std::cout << "Failed assert in template/function " << myTemplateName << " line 35. " <<  "Followed trace of components: " << ctx->getTrace(myId) << std::endl;
assert(Fr_isTrue(&expaux[0]));
for (uint i = 0; i < 0; i++){
uint index_subc = ctx->componentMemory[ctx_index].subcomponents[i];
if (index_subc != 0)release_memory_component(ctx,index_subc);
}
}

void OneToNine_3_create(uint soffset,uint coffset,Circom_CalcWit* ctx,std::string componentName,uint componentFather){
ctx->componentMemory[coffset].templateId = 3;
ctx->componentMemory[coffset].templateName = "OneToNine";
ctx->componentMemory[coffset].signalStart = soffset;
ctx->componentMemory[coffset].inputCounter = 1;
ctx->componentMemory[coffset].componentName = componentName;
ctx->componentMemory[coffset].idFather = componentFather;
ctx->componentMemory[coffset].subcomponents = new uint[2]{0};
}

void OneToNine_3_run(uint ctx_index,Circom_CalcWit* ctx){
FrElement* signalValues = ctx->signalValues;
u64 mySignalStart = ctx->componentMemory[ctx_index].signalStart;
std::string myTemplateName = ctx->componentMemory[ctx_index].templateName;
std::string myComponentName = ctx->componentMemory[ctx_index].componentName;
u64 myFather = ctx->componentMemory[ctx_index].idFather;
u64 myId = ctx_index;
u32* mySubcomponents = ctx->componentMemory[ctx_index].subcomponents;
bool* mySubcomponentsParallel = ctx->componentMemory[ctx_index].subcomponentsParallel;
FrElement* circuitConstants = ctx->circuitConstants;
std::string* listOfTemplateMessages = ctx->listOfTemplateMessages;
FrElement expaux[3];
FrElement lvar[0];
uint sub_component_aux;
uint index_multiple_eq;
{
uint aux_create = 0;
int aux_cmp_num = 0+ctx_index+1;
uint csoffset = mySignalStart+1;
for (uint i = 0; i < 1; i++) {
std::string new_cmp_name = "lowerBound";
Bits4_2_create(csoffset,aux_cmp_num,ctx,new_cmp_name,myId);
mySubcomponents[aux_create+i] = aux_cmp_num;
csoffset += 5 ;
aux_cmp_num += 1;
}
}
{
uint aux_create = 1;
int aux_cmp_num = 1+ctx_index+1;
uint csoffset = mySignalStart+6;
for (uint i = 0; i < 1; i++) {
std::string new_cmp_name = "upperBound";
Bits4_2_create(csoffset,aux_cmp_num,ctx,new_cmp_name,myId);
mySubcomponents[aux_create+i] = aux_cmp_num;
csoffset += 5 ;
aux_cmp_num += 1;
}
}
{
uint cmp_index_ref = 0;
{
PFrElement aux_dest = &ctx->signalValues[ctx->componentMemory[mySubcomponents[cmp_index_ref]].signalStart + 0];
// load src
Fr_sub(&expaux[0],&signalValues[mySignalStart + 0],&circuitConstants[0]); // line circom 43
// end load src
Fr_copy(aux_dest,&expaux[0]);
}
// need to run sub component
assert(!(ctx->componentMemory[mySubcomponents[cmp_index_ref]].inputCounter -= 1));
Bits4_2_run(mySubcomponents[cmp_index_ref],ctx);
}
{
uint cmp_index_ref = 1;
{
PFrElement aux_dest = &ctx->signalValues[ctx->componentMemory[mySubcomponents[cmp_index_ref]].signalStart + 0];
// load src
Fr_add(&expaux[0],&signalValues[mySignalStart + 0],&circuitConstants[5]); // line circom 44
// end load src
Fr_copy(aux_dest,&expaux[0]);
}
// need to run sub component
assert(!(ctx->componentMemory[mySubcomponents[cmp_index_ref]].inputCounter -= 1));
Bits4_2_run(mySubcomponents[cmp_index_ref],ctx);
}
for (uint i = 0; i < 2; i++){
uint index_subc = ctx->componentMemory[ctx_index].subcomponents[i];
if (index_subc != 0)release_memory_component(ctx,index_subc);
}
}

void Sudoku_4_create(uint soffset,uint coffset,Circom_CalcWit* ctx,std::string componentName,uint componentFather){
ctx->componentMemory[coffset].templateId = 4;
ctx->componentMemory[coffset].templateName = "Sudoku";
ctx->componentMemory[coffset].signalStart = soffset;
ctx->componentMemory[coffset].inputCounter = 162;
ctx->componentMemory[coffset].componentName = componentName;
ctx->componentMemory[coffset].idFather = componentFather;
ctx->componentMemory[coffset].subcomponents = new uint[90]{0};
}

void Sudoku_4_run(uint ctx_index,Circom_CalcWit* ctx){
FrElement* signalValues = ctx->signalValues;
u64 mySignalStart = ctx->componentMemory[ctx_index].signalStart;
std::string myTemplateName = ctx->componentMemory[ctx_index].templateName;
std::string myComponentName = ctx->componentMemory[ctx_index].componentName;
u64 myFather = ctx->componentMemory[ctx_index].idFather;
u64 myId = ctx_index;
u32* mySubcomponents = ctx->componentMemory[ctx_index].subcomponents;
bool* mySubcomponentsParallel = ctx->componentMemory[ctx_index].subcomponentsParallel;
FrElement* circuitConstants = ctx->circuitConstants;
std::string* listOfTemplateMessages = ctx->listOfTemplateMessages;
FrElement expaux[6];
FrElement lvar[3];
uint sub_component_aux;
uint index_multiple_eq;
{
PFrElement aux_dest = &lvar[0];
// load src
// end load src
Fr_copy(aux_dest,&circuitConstants[1]);
}
{
uint aux_create = 0;
int aux_cmp_num = 0+ctx_index+1;
uint csoffset = mySignalStart+162;
uint aux_dimensions[1] = {9};
for (uint i = 0; i < 9; i++) {
std::string new_cmp_name = "distinct"+ctx->generate_position_array(aux_dimensions, 1, i);
Distinct_1_create(csoffset,aux_cmp_num,ctx,new_cmp_name,myId);
mySubcomponents[aux_create+i] = aux_cmp_num;
csoffset += 117 ;
aux_cmp_num += 37;
}
}
{
uint aux_create = 9;
int aux_cmp_num = 333+ctx_index+1;
uint csoffset = mySignalStart+1215;
uint aux_dimensions[2] = {9,9};
for (uint i = 0; i < 81; i++) {
std::string new_cmp_name = "inRange"+ctx->generate_position_array(aux_dimensions, 2, i);
OneToNine_3_create(csoffset,aux_cmp_num,ctx,new_cmp_name,myId);
mySubcomponents[aux_create+i] = aux_cmp_num;
csoffset += 11 ;
aux_cmp_num += 3;
}
}
{
PFrElement aux_dest = &lvar[1];
// load src
// end load src
Fr_copy(aux_dest,&circuitConstants[2]);
}
Fr_lt(&expaux[0],&lvar[1],&circuitConstants[1]); // line circom 56
while(Fr_isTrue(&expaux[0])){
{
PFrElement aux_dest = &lvar[2];
// load src
// end load src
Fr_copy(aux_dest,&circuitConstants[2]);
}
Fr_lt(&expaux[0],&lvar[2],&circuitConstants[1]); // line circom 57
while(Fr_isTrue(&expaux[0])){
Fr_sub(&expaux[3],&signalValues[mySignalStart + (((9 * Fr_toInt(&lvar[1])) + (1 * Fr_toInt(&lvar[2]))) + 0)],&signalValues[mySignalStart + (((9 * Fr_toInt(&lvar[1])) + (1 * Fr_toInt(&lvar[2]))) + 81)]); // line circom 59
Fr_mul(&expaux[1],&signalValues[mySignalStart + (((9 * Fr_toInt(&lvar[1])) + (1 * Fr_toInt(&lvar[2]))) + 0)],&expaux[3]); // line circom 59
Fr_eq(&expaux[0],&expaux[1],&circuitConstants[2]); // line circom 59
if (!Fr_isTrue(&expaux[0])) std::cout << "Failed assert in template/function " << myTemplateName << " line 59. " <<  "Followed trace of components: " << ctx->getTrace(myId) << std::endl;
assert(Fr_isTrue(&expaux[0]));
{
PFrElement aux_dest = &lvar[2];
// load src
Fr_add(&expaux[0],&lvar[2],&circuitConstants[0]); // line circom 57
// end load src
Fr_copy(aux_dest,&expaux[0]);
}
Fr_lt(&expaux[0],&lvar[2],&circuitConstants[1]); // line circom 57
}
{
PFrElement aux_dest = &lvar[1];
// load src
Fr_add(&expaux[0],&lvar[1],&circuitConstants[0]); // line circom 56
// end load src
Fr_copy(aux_dest,&expaux[0]);
}
Fr_lt(&expaux[0],&lvar[1],&circuitConstants[1]); // line circom 56
}
{
PFrElement aux_dest = &lvar[1];
// load src
// end load src
Fr_copy(aux_dest,&circuitConstants[2]);
}
Fr_lt(&expaux[0],&lvar[1],&circuitConstants[1]); // line circom 63
while(Fr_isTrue(&expaux[0])){
{
PFrElement aux_dest = &lvar[2];
// load src
// end load src
Fr_copy(aux_dest,&circuitConstants[2]);
}
Fr_lt(&expaux[0],&lvar[2],&circuitConstants[1]); // line circom 64
while(Fr_isTrue(&expaux[0])){
Fr_eq(&expaux[0],&lvar[1],&circuitConstants[2]); // line circom 65
if(Fr_isTrue(&expaux[0])){

}
{
uint cmp_index_ref = (((9 * Fr_toInt(&lvar[1])) + (1 * Fr_toInt(&lvar[2]))) + 9);
{
PFrElement aux_dest = &ctx->signalValues[ctx->componentMemory[mySubcomponents[cmp_index_ref]].signalStart + 0];
// load src
// end load src
Fr_copy(aux_dest,&signalValues[mySignalStart + (((9 * Fr_toInt(&lvar[1])) + (1 * Fr_toInt(&lvar[2]))) + 81)]);
}
// run sub component if needed
if(!(ctx->componentMemory[mySubcomponents[cmp_index_ref]].inputCounter -= 1)){
OneToNine_3_run(mySubcomponents[cmp_index_ref],ctx);

}
}
{
uint cmp_index_ref = ((1 * Fr_toInt(&lvar[2])) + 0);
{
PFrElement aux_dest = &ctx->signalValues[ctx->componentMemory[mySubcomponents[cmp_index_ref]].signalStart + ((1 * Fr_toInt(&lvar[1])) + 0)];
// load src
// end load src
Fr_copy(aux_dest,&signalValues[mySignalStart + (((9 * Fr_toInt(&lvar[1])) + (1 * Fr_toInt(&lvar[2]))) + 81)]);
}
// run sub component if needed
if(!(ctx->componentMemory[mySubcomponents[cmp_index_ref]].inputCounter -= 1)){
Distinct_1_run(mySubcomponents[cmp_index_ref],ctx);

}
}
{
PFrElement aux_dest = &lvar[2];
// load src
Fr_add(&expaux[0],&lvar[2],&circuitConstants[0]); // line circom 64
// end load src
Fr_copy(aux_dest,&expaux[0]);
}
Fr_lt(&expaux[0],&lvar[2],&circuitConstants[1]); // line circom 64
}
{
PFrElement aux_dest = &lvar[1];
// load src
Fr_add(&expaux[0],&lvar[1],&circuitConstants[0]); // line circom 63
// end load src
Fr_copy(aux_dest,&expaux[0]);
}
Fr_lt(&expaux[0],&lvar[1],&circuitConstants[1]); // line circom 63
}
for (uint i = 0; i < 90; i++){
uint index_subc = ctx->componentMemory[ctx_index].subcomponents[i];
if (index_subc != 0)release_memory_component(ctx,index_subc);
}
}

void run(Circom_CalcWit* ctx){
Sudoku_4_create(1,0,ctx,"main",0);
Sudoku_4_run(0,ctx);
}

