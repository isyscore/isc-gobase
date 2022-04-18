#ifndef __GOJVM_H__
#define __GOJVM_H__

#include <stdio.h>
#include <stdlib.h>
#include <jni.h>
#include <string.h>
#include <stdbool.h>
#include "gojvm_wrap.h"

JavaVM* createJvm(char* classPath, char* xms, char* xmx, char* xmn, char* xss);
int destroyJvm(JavaVM* jvm);
JNIEnv* attachJvm(JavaVM* jvm);
void detachJvm(JavaVM* jvm);
jclass findClass(JNIEnv* env, char* className);
jclass getObjectClass(JNIEnv* env, jobject obj);
jobject newJavaObject(JNIEnv* env, jclass clazz);
void freeJavaClassRef(JNIEnv* env, jclass clz);
void freeJavaObject(JNIEnv* env, jobject obj);

void callStaticVoidMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);
jobject callStaticObjectMethod(JNIEnv *env, jclass clazz, char *methodName, char *sig, int len, char **types, void **args);
char* callStaticStringMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);
int callStaticIntMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);
long callStaticLongMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);
short callStaticShortMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);
unsigned char callStaticByteMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);
float callStaticFloatMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);
double callStaticDoubleMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);
int callStaticBooleanMethod(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args);

jobject getStaticObject(JNIEnv* env, jclass clazz, char* fieldName, char* sig);
void setStaticObject(JNIEnv* env, jclass clazz, char* fieldName, char* sig, jobject obj);
char* getStaticString(JNIEnv* env, jclass clazz, char* fieldName);
void setStaticString(JNIEnv* env, jclass clazz, char* fieldName, char* value);
int getStaticInt(JNIEnv* env, jclass clazz, char* fieldName);
void setStaticInt(JNIEnv* env, jclass clazz, char* fieldName, int value);
long getStaticLong(JNIEnv *env, jclass clazz, char *fieldName);
void setStaticLong(JNIEnv *env, jclass clazz, char *fieldName, long value);
short getStaticShort(JNIEnv *env, jclass clazz, char *fieldName);
void setStaticShort(JNIEnv *env, jclass clazz, char *fieldName, short value);
unsigned char getStaticByte(JNIEnv *env, jclass clazz, char *fieldName);
void setStaticByte(JNIEnv *env, jclass clazz, char *fieldName, unsigned char value);
float getStaticFloat(JNIEnv *env, jclass clazz, char *fieldName);
void setStaticFloat(JNIEnv *env, jclass clazz, char *fieldName, float value);
double getStaticDouble(JNIEnv *env, jclass clazz, char *fieldName);
void setStaticDouble(JNIEnv *env, jclass clazz, char *fieldName, double value);
int getStaticBoolean(JNIEnv *env, jclass clazz, char *fieldName);
void setStaticBoolean(JNIEnv *env, jclass clazz, char *fieldName, int value);

void callVoidMethod(JNIEnv *env, jclass clazz, jobject obj, char *methodName, char *sig, int len, char **types, void **args);
jobject callObjectMethod(JNIEnv *env, jclass clazz, jobject obj, char *methodName, char *sig, int len, char **types, void **args);
char *callStringMethod(JNIEnv *env, jclass clazz, jobject obj, char *methodName, char *sig, int len, char **types, void **args);
int callIntMethod(JNIEnv *env, jclass clazz, jobject obj, char *methodName, char *sig, int len, char **types, void **args);
long callLongMethod(JNIEnv *env, jclass clazz, jobject obj, char *methodName, char *sig, int len, char **types, void **args);
short callShortMethod(JNIEnv *env, jclass clazz, jobject obj, char *methodName, char *sig, int len, char **types, void **args);
unsigned char callByteMethod(JNIEnv *env, jclass clazz, jobject obj, char *methodName, char *sig, int len, char **types, void **args);
float callFloatMethod(JNIEnv *env, jclass clazz, jobject obj, char *methodName, char *sig, int len, char **types, void **args);
double callDoubleMethod(JNIEnv *env, jclass clazz, jobject obj, char *methodName, char *sig, int len, char **types, void **args);
int callBooleanMethod(JNIEnv *env, jclass clazz, jobject obj, char *methodName, char *sig, int len, char **types, void **args);

jobject getObjectObject(JNIEnv *env, jclass clazz, jobject obj, char *fieldName, char *sig);
void setObjectObject(JNIEnv *env, jclass clazz, jobject obj, char *fieldName, char *sig, jobject value);
char* getObjectString(JNIEnv* env, jclass clazz, jobject obj, char* fieldName);
void setObjectString(JNIEnv* env, jclass clazz, jobject obj, char* fieldName, char* value);
int getObjectInt(JNIEnv *env, jclass clazz, jobject obj, char *fieldName);
void setObjectInt(JNIEnv *env, jclass clazz, jobject obj, char *fieldName, int value);
long getObjectLong(JNIEnv *env, jclass clazz, jobject obj, char *fieldName);
void setObjectLong(JNIEnv *env, jclass clazz, jobject obj, char *fieldName, long value);
short getObjectShort(JNIEnv *env, jclass clazz, jobject obj, char *fieldName);
void setObjectShort(JNIEnv *env, jclass clazz, jobject obj, char *fieldName, short value);
unsigned char getObjectByte(JNIEnv *env, jclass clazz, jobject obj, char *fieldName);
void setObjectByte(JNIEnv *env, jclass clazz, jobject obj, char *fieldName, unsigned char value);
float getObjectFloat(JNIEnv *env, jclass clazz, jobject obj, char *fieldName);
void setObjectFloat(JNIEnv *env, jclass clazz, jobject obj, char *fieldName, float value);
double getObjectDouble(JNIEnv *env, jclass clazz, jobject obj, char *fieldName);
void setObjectDouble(JNIEnv *env, jclass clazz, jobject obj, char *fieldName, double value);
int getObjectBoolean(JNIEnv *env, jclass clazz, jobject obj, char *fieldName);
void setObjectBoolean(JNIEnv *env, jclass clazz, jobject obj, char *fieldName, int value);

#endif // __GOJVM_H__