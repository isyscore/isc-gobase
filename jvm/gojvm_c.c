#include "gojvm_c.h"

bool hasPrefix(const char *str, const char *sub) {
    return strncmp(str, sub, strlen(sub)) == 0;
}

char* getSubType(char* type) {
    size_t len = strlen(type);
    char* s = (char*)malloc(len);
    int inType = 0;
    int iidx = 0;
    for (int i = 0; i < len; i++) {
        if (type[i] == '<') {
            inType = 1;
            continue;
        }
        if (type[i] == '>') break;
        if (inType) s[iidx++] = type[i];
    }
    return s;
}

char* getRealSig(char* sig) {
    size_t len = strlen(sig);
    char* s = (char*) malloc(len);
    int inType = 0;
    int iidx = 0;
    for (int i = 0; i < len; i++) {
        if (sig[i] == '<') {
            inType = 1;
            continue;
        }
        if (sig[i] == '>') {
            inType = 0;
            continue;
        }
        if (!inType) {
            s[iidx++] = sig[i];
        }
    }
    return s;
}

jvalue* makeParams(JNIEnv* env, int len, char** types, void** args) {
    jvalue *v = malloc(sizeof(jvalue) * len);
    for (int i = 0; i < len; i++) {
        if (strcmp(types[i], "Ljava/lang/String;") == 0) {
            v[i].l = (*env)->NewStringUTF(env, (char*)args[i]);
        } else if (strcmp(types[i], "I") == 0) {
            v[i].i = *((int*)args[i]);
        } else if (strcmp(types[i], "J") == 0) {
            v[i].j = *((long*)args[i]);
        } else if (strcmp(types[i], "F") == 0) {
            v[i].f = *((float*)args[i]);
        } else if (strcmp(types[i], "D") == 0) {
            v[i].d = *((double*)args[i]);
        } else if (strcmp(types[i], "B") == 0) {
            v[i].b = *((unsigned char*)args[i]);
        } else if (strcmp(types[i], "S") == 0) {
            v[i].s = *((short*)args[i]);
        } else if (strcmp(types[i], "Z") == 0) {
            int bi = *((int*)args[i]);
            v[i].z = bi == 0 ? JNI_FALSE : JNI_TRUE;
        }
    }
    return v;
}

void freeParams(JNIEnv* env, int len, char** types, jvalue* v) {
    for (int i = 0; i < len; i++) {
        if (strcmp(types[i], "Ljava/lang/String;") == 0) {
            (*env)->DeleteLocalRef(env, v[i].l);
        }
    }
    free(v);
}

_GO_EXPORT JavaVM* createJvm(char* classPath, char* xms, char* xmx, char* xmn, char* xss) {
	JavaVM* jvm;
	JNIEnv* env;
	JavaVMInitArgs vm_args;
	JavaVMOption options[5];

	options[0].optionString = (char*)malloc(strlen("-Djava.class.path=") + strlen(classPath) + 1);
	sprintf(options[0].optionString, "-Djava.class.path=%s", classPath);
	options[1].optionString = (char*)malloc(strlen("-Xms") + strlen(xms) + 1);
	sprintf(options[1].optionString, "-Xms%s", xms);
	options[2].optionString = (char*)malloc(strlen("-Xmx") + strlen(xmx) + 1);
	sprintf(options[2].optionString, "-Xmx%s", xmx);
	options[3].optionString = (char*)malloc(strlen("-Xmn") + strlen(xmn) + 1);
	sprintf(options[3].optionString, "-Xmn%s", xmn);
	options[4].optionString = (char*)malloc(strlen("-Xss") + strlen(xss) + 1);
	sprintf(options[4].optionString, "-Xss%s", xss);

	vm_args.version = JNI_VERSION_1_8;
	vm_args.nOptions = 5;
	vm_args.options = options;
	vm_args.ignoreUnrecognized = JNI_FALSE;

	jint res = JNI_CreateJavaVM(&jvm, (void**)&env, &vm_args);
	if (res < 0) {
		printf("create jvm failed\n");
		return NULL;
	}
	(*jvm)->DetachCurrentThread(jvm);
	return jvm;
}

_GO_EXPORT int destroyJvm(JavaVM* jvm) {
	jint res = (*jvm)->DestroyJavaVM(jvm);
	if (res < 0) {
		printf("destroy jvm failed\n");
		return 1;
	}
	return 0;
}


_GO_EXPORT JNIEnv* attachJvm(JavaVM* jvm) {
	JNIEnv* env;
	jint res = (*jvm)->AttachCurrentThread(jvm, (void**)&env, NULL);
	if (res < 0) {
		printf("attach jvm failed\n");
		return NULL;
	}
	return env;
}

_GO_EXPORT void detachJvm(JavaVM* jvm) {
	(*jvm)->DetachCurrentThread(jvm);
}

_GO_EXPORT jclass findClass(JNIEnv* env, char* className) {
	jclass cls = (*env)->FindClass(env, className);
	if (cls == NULL) {
		printf("find class failed\n");
		return NULL;
	}
	return cls;
}

_GO_EXPORT jclass getObjectClass(JNIEnv* env, jobject obj) {
    return (*env)->GetObjectClass(env, obj);
}

_GO_EXPORT jobject newJavaObject(JNIEnv* env, jclass clazz) {
    jmethodID m = (*env)->GetMethodID(env, clazz, "<init>", "()V");
    jobject jret = (*env)->NewObject(env, clazz, m);
    return jret;
}

_GO_EXPORT void freeJavaObject(JNIEnv* env, jobject obj) {
    (*env)->DeleteLocalRef(env, obj);
}

_GO_EXPORT void freeJavaClassRef(JNIEnv* env, jclass clz) {
    (*env)->DeleteLocalRef(env, clz);
}

_GO_EXPORT WRAP_STATIC_VOID_METHOD(Void)
_GO_EXPORT WRAP_STATIC_METHOD(Object, jobject, jobject)
_GO_EXPORT WRAP_STATIC_STRING_METHOD(String)
_GO_EXPORT WRAP_STATIC_METHOD(Int, int, jint)
_GO_EXPORT WRAP_STATIC_METHOD(Long, long, jlong)
_GO_EXPORT WRAP_STATIC_METHOD(Short, short, jshort)
_GO_EXPORT WRAP_STATIC_METHOD(Byte, unsigned char, jbyte)
_GO_EXPORT WRAP_STATIC_METHOD(Float, float, jfloat)
_GO_EXPORT WRAP_STATIC_METHOD(Double, double, jdouble)
_GO_EXPORT WRAP_STATIC_METHOD(Boolean, int, jboolean)

_GO_EXPORT WRAP_STATIC_FIELD_GET(Object, jobject, jobject)
_GO_EXPORT WRAP_STATIC_FIELD_SET(Object, jobject)
_GO_EXPORT WRAP_STATIC_FIELD_GET_STRING(String)
_GO_EXPORT WRAP_STATIC_FIELD_SET_STRING(String)
_GO_EXPORT WRAP_STATIC_FIELD_GET_SIG(Int, int, jint, "I")
_GO_EXPORT WRAP_STATIC_FIELD_SET_SIG(Int, int, jint, "I")
_GO_EXPORT WRAP_STATIC_FIELD_GET_SIG(Long, long, jlong, "J")
_GO_EXPORT WRAP_STATIC_FIELD_SET_SIG(Long, long, jlong, "J")
_GO_EXPORT WRAP_STATIC_FIELD_GET_SIG(Short, short, jshort , "S")
_GO_EXPORT WRAP_STATIC_FIELD_SET_SIG(Short, short, jshort, "S")
_GO_EXPORT WRAP_STATIC_FIELD_GET_SIG(Byte, unsigned char, jbyte , "B")
_GO_EXPORT WRAP_STATIC_FIELD_SET_SIG(Byte, unsigned char, jbyte, "B")
_GO_EXPORT WRAP_STATIC_FIELD_GET_SIG(Float, float, jfloat, "F")
_GO_EXPORT WRAP_STATIC_FIELD_SET_SIG(Float, float, jfloat, "F")
_GO_EXPORT WRAP_STATIC_FIELD_GET_SIG(Double, double, jdouble, "D")
_GO_EXPORT WRAP_STATIC_FIELD_SET_SIG(Double, double, jdouble, "D")
_GO_EXPORT WRAP_STATIC_FIELD_GET_SIG(Boolean, int, jboolean , "Z")
_GO_EXPORT WRAP_STATIC_FIELD_SET_SIG(Boolean, int, jboolean, "Z")

_GO_EXPORT WRAP_VOID_METHOD(Void)
_GO_EXPORT WRAP_METHOD(Object, jobject, jobject)
_GO_EXPORT WRAP_STRING_METHOD(String)
_GO_EXPORT WRAP_METHOD(Int, int, jint)
_GO_EXPORT WRAP_METHOD(Long, long, jlong)
_GO_EXPORT WRAP_METHOD(Short, short, jshort)
_GO_EXPORT WRAP_METHOD(Byte, unsigned char, jbyte)
_GO_EXPORT WRAP_METHOD(Float, float, jfloat)
_GO_EXPORT WRAP_METHOD(Double, double, jdouble)
_GO_EXPORT WRAP_METHOD(Boolean, int, jboolean)

_GO_EXPORT WRAP_FIELD_GET(Object, jobject, jobject)
_GO_EXPORT WRAP_FIELD_SET(Object, jobject)
_GO_EXPORT WRAP_FIELD_GET_STRING(String)
_GO_EXPORT WRAP_FIELD_SET_STRING(String)
_GO_EXPORT WRAP_FIELD_GET_SIG(Int, int, jint, "I")
_GO_EXPORT WRAP_FIELD_SET_SIG(Int, int, jint, "I")
_GO_EXPORT WRAP_FIELD_GET_SIG(Long, long, jlong, "J")
_GO_EXPORT WRAP_FIELD_SET_SIG(Long, long, jlong, "J")
_GO_EXPORT WRAP_FIELD_GET_SIG(Short, short, jshort , "S")
_GO_EXPORT WRAP_FIELD_SET_SIG(Short, short, jshort, "S")
_GO_EXPORT WRAP_FIELD_GET_SIG(Byte, unsigned char, jbyte , "B")
_GO_EXPORT WRAP_FIELD_SET_SIG(Byte, unsigned char, jbyte, "B")
_GO_EXPORT WRAP_FIELD_GET_SIG(Float, float, jfloat, "F")
_GO_EXPORT WRAP_FIELD_SET_SIG(Float, float, jfloat, "F")
_GO_EXPORT WRAP_FIELD_GET_SIG(Double, double, jdouble, "D")
_GO_EXPORT WRAP_FIELD_SET_SIG(Double, double, jdouble, "D")
_GO_EXPORT WRAP_FIELD_GET_SIG(Boolean, int, jboolean , "Z")
_GO_EXPORT WRAP_FIELD_SET_SIG(Boolean, int, jboolean, "Z")
