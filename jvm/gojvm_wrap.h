//
// Created by rarnu on 2022/3/29.
//

#ifndef GOJVM_GOJVM_WRAP_H
#define GOJVM_GOJVM_WRAP_H

#define _GO_EXPORT __attribute__((__visibility__("default")))

#define WRAP_RETURN_STRING \
    const char* str = (*env)->GetStringUTFChars(env, jret, NULL); \
    (*env)->DeleteLocalRef(env, jret);                            \
    return (char*)str;

#define WRAP_STATIC_VOID(JNINAME) \
    jmethodID m = (*env)->GetStaticMethodID(env, clazz, methodName, sig);   \
    jvalue *v = makeParams(env, len, types, args);                          \
    (*env)->JNINAME(env, clazz, m, v);                                      \
    freeParams(env, len, types, v);

#define WRAP_VOID(JNINAME) \
    jmethodID m = (*env)->GetMethodID(env, clazz, methodName, sig);   \
    jvalue *v = makeParams(env, len, types, args);                    \
    (*env)->JNINAME(env, obj, m, v);                                  \
    freeParams(env, len, types, v);

#define WRAP_STATIC_RET(JNINAME, JNIRET) \
    jmethodID m = (*env)->GetStaticMethodID(env, clazz, methodName, sig);   \
    jvalue *v = makeParams(env, len, types, args);                          \
    JNIRET jret = (*env)->JNINAME(env, clazz, m, v);                        \
    freeParams(env, len, types, v);

#define WRAP_RET(JNINAME, JNIRET) \
    jmethodID m = (*env)->GetMethodID(env, clazz, methodName, sig);   \
    jvalue *v = makeParams(env, len, types, args);                    \
    JNIRET jret = (*env)->JNINAME(env, obj, m, v);                    \
    freeParams(env, len, types, v);

#define WRAP_STATIC_VOID_METHOD(NAME) \
    void callStatic##NAME##Method(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args) { \
        WRAP_STATIC_VOID(CallStatic##NAME##MethodA)                                                                             \
    }

#define WRAP_VOID_METHOD(NAME) \
    void call##NAME##Method(JNIEnv* env, jclass clazz, jobject obj, char* methodName, char* sig, int len, char** types, void** args) { \
        WRAP_VOID(Call##NAME##MethodA)                                                                                                 \
    }

#define WRAP_STATIC_METHOD(NAME, RET, JNIRET) \
    RET callStatic##NAME##Method(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args) {  \
        WRAP_STATIC_RET(CallStatic##NAME##MethodA, JNIRET)                                                                      \
        return (RET)jret;                                                                                                       \
    }

#define WRAP_METHOD(NAME, RET, JNIRET) \
    RET call##NAME##Method(JNIEnv* env, jclass clazz, jobject obj, char* methodName, char* sig, int len, char** types, void** args) {  \
        WRAP_RET(Call##NAME##MethodA, JNIRET)                                                                                          \
        return (RET)jret;                                                                                                              \
    }

#define WRAP_STATIC_STRING_METHOD(NAME) \
    char* callStatic##NAME##Method(JNIEnv* env, jclass clazz, char* methodName, char* sig, int len, char** types, void** args) {  \
        WRAP_STATIC_RET(CallStaticObjectMethodA, jobject)                                                                         \
        WRAP_RETURN_STRING                                                                                                        \
    }

#define WRAP_STRING_METHOD(NAME) \
    char* call##NAME##Method(JNIEnv* env, jclass clazz, jobject obj, char* methodName, char* sig, int len, char** types, void** args) {  \
        WRAP_RET(CallObjectMethodA, jobject)                                                                                             \
        WRAP_RETURN_STRING                                                                                                               \
    }

#define WRAP_STATIC_FIELD_GET(NAME, RET, JNIRET) \
    RET getStatic##NAME(JNIEnv* env, jclass clazz, char* fieldName, char* sig) { \
        jfieldID f = (*env)->GetStaticFieldID(env, clazz, fieldName, sig);       \
        JNIRET jret = (*env)->GetStatic##NAME##Field(env, clazz, f);             \
        return (RET) jret;                                                      \
    }

#define WRAP_STATIC_FIELD_SET(NAME, JNITYPE) \
    void setStatic##NAME(JNIEnv* env, jclass clazz, char* fieldName, char* sig, JNITYPE obj) { \
        jfieldID f = (*env)->GetStaticFieldID(env, clazz, fieldName, sig);                     \
        (*env)->SetStatic##NAME##Field(env, clazz, f, obj);                                    \
    }

#define WRAP_STATIC_FIELD_GET_STRING(NAME) \
    char* getStatic##NAME(JNIEnv* env, jclass clazz, char* fieldName) {                 \
        jstring jret = getStaticObject(env, clazz, fieldName, "Ljava/lang/String;");    \
        WRAP_RETURN_STRING                                                              \
    }

#define WRAP_STATIC_FIELD_SET_STRING(NAME) \
    void setStatic##NAME(JNIEnv* env, jclass clazz, char* fieldName, char* value) { \
        jstring jstr = (*env)->NewStringUTF(env, value);                            \
        setStaticObject(env, clazz, fieldName, "Ljava/lang/String;", jstr);         \
        (*env)->DeleteLocalRef(env, jstr);                                          \
    }

#define WRAP_STATIC_FIELD_GET_SIG(NAME, RET, JNIRET, SIG) \
    RET getStatic##NAME(JNIEnv* env, jclass clazz, char* fieldName) {      \
        jfieldID f = (*env)->GetStaticFieldID(env, clazz, fieldName, SIG); \
        JNIRET jret = (*env)->GetStatic##NAME##Field(env, clazz, f);       \
        return (RET)jret;                                                  \
    }

#define WRAP_STATIC_FIELD_SET_SIG(NAME, TYPE, JNITYPE, SIG) \
    void setStatic##NAME(JNIEnv* env, jclass clazz, char* fieldName, TYPE value) { \
        jfieldID f = (*env)->GetStaticFieldID(env, clazz, fieldName, SIG);         \
        (*env)->SetStatic##NAME##Field(env, clazz, f, (JNITYPE)value);             \
    }

#define WRAP_FIELD_GET(NAME, RET, JNIRET) \
    RET getObject##NAME(JNIEnv* env, jclass clazz, jobject obj, char* fieldName, char* sig) { \
        jfieldID f = (*env)->GetFieldID(env, clazz, fieldName, sig);                          \
        JNIRET jret = (*env)->Get##NAME##Field(env, obj, f);                                  \
        return (RET) jret;                                                                    \
    }

#define WRAP_FIELD_SET(NAME, JNITYPE) \
    void setObject##NAME(JNIEnv* env, jclass clazz, jobject obj, char* fieldName, char* sig, JNITYPE value) { \
        jfieldID f = (*env)->GetFieldID(env, clazz, fieldName, sig);                                          \
        (*env)->Set##NAME##Field(env, obj, f, value);                                                         \
    }

#define WRAP_FIELD_GET_STRING(NAME) \
    char* getObject##NAME(JNIEnv* env, jclass clazz, jobject obj, char* fieldName) {      \
        jstring jret = getObjectObject(env, clazz, obj, fieldName, "Ljava/lang/String;"); \
        WRAP_RETURN_STRING                                                                \
    }

#define WRAP_FIELD_SET_STRING(NAME) \
    void setObject##NAME(JNIEnv* env, jclass clazz, jobject obj, char* fieldName, char* value) { \
        jstring jstr = (*env)->NewStringUTF(env, value);                                         \
        setObjectObject(env, clazz, obj, fieldName, "Ljava/lang/String;", jstr);                 \
        (*env)->DeleteLocalRef(env, jstr);                                                       \
    }

#define WRAP_FIELD_GET_SIG(NAME, RET, JNIRET, SIG) \
    RET getObject##NAME(JNIEnv* env, jclass clazz, jobject obj, char* fieldName) { \
        jfieldID f = (*env)->GetFieldID(env, clazz, fieldName, SIG);               \
        JNIRET jret = (*env)->Get##NAME##Field(env, obj, f);                       \
        return (RET)jret;                                                          \
    }

#define WRAP_FIELD_SET_SIG(NAME, TYPE, JNITYPE, SIG) \
    void setObject##NAME(JNIEnv* env, jclass clazz, jobject obj, char* fieldName, TYPE value) { \
        jfieldID f = (*env)->GetFieldID(env, clazz, fieldName, SIG);                            \
        (*env)->Set##NAME##Field(env, obj, f, (JNITYPE)value);                                  \
    }

#endif //GOJVM_GOJVM_WRAP_H
