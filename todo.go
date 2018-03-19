package unqlitego

/* TODO: implement

// Database Engine Handle
int unqlite_config(unqlite *pDb,int nOp,...);

// Key/Value (KV) Store Interfaces
int unqlite_kv_fetch_callback(unqlite *pDb,const void *pKey,
	                    int nKeyLen,int (*xConsumer)(const void *,unsigned int,void *),void *pUserData);
int unqlite_kv_config(unqlite *pDb,int iOp,...);

//  Cursor Iterator Interfaces
int unqlite_kv_cursor_key_callback(unqlite_kv_cursor *pCursor,int (*xConsumer)(const void *,unsigned int,void *),void *pUserData);
int unqlite_kv_cursor_data_callback(unqlite_kv_cursor *pCursor,int (*xConsumer)(const void *,unsigned int,void *),void *pUserData);

// Utility interfaces
int unqlite_util_load_mmaped_file(const char *zFile,void **ppMap,unqlite_int64 *pFileSize);
int unqlite_util_release_mmaped_file(void *pMap,unqlite_int64 iFileSize);

// Global Library Management Interfaces
int unqlite_lib_config(int nConfigOp,...);
*/
