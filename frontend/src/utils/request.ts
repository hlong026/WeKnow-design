// src/utils/request.js
import axios from "axios";
import { generateRandomString } from "./index";

// API基础URL
const BASE_URL = import.meta.env.VITE_IS_DOCKER ? "" : "http://localhost:8080";

// 创建Axios实例
const instance = axios.create({
  baseURL: BASE_URL,
  timeout: 30000,
  headers: {
    "Content-Type": "application/json",
    "X-Request-ID": `${generateRandomString(12)}`,
  },
});

// 单机版：简化的请求拦截器
// 在 disable_auth 模式下，后端会自动处理认证，前端不需要发送任何认证信息
instance.interceptors.request.use(
  (config) => {
    config.headers["X-Request-ID"] = `${generateRandomString(12)}`;
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 单机版：简化的响应拦截器
instance.interceptors.response.use(
  (response) => {
    const { status, data } = response;
    if (status === 200 || status === 201) {
      return data;
    } else {
      return Promise.reject(data);
    }
  },
  async (error: any) => {
    if (!error.response) {
      return Promise.reject({ message: "网络错误，请检查您的网络连接" });
    }
    
    // 处理 Nginx 413 Request Entity Too Large
    if (error.response.status === 413) {
      return Promise.reject({ 
        status: 413, 
        message: '文件大小超过限制，请上传较小的文件',
        success: false
      });
    }

    const { status, data } = error.response;
    const errorMessage = typeof data === 'object' && data?.error?.message 
      ? data.error.message 
      : (typeof data === 'object' ? data?.message : data);
    return Promise.reject({ 
      status, 
      message: errorMessage,
      ...(typeof data === 'object' ? data : {}) 
    });
  }
);

export function get(url: string) {
  return instance.get(url);
}

export async function getDown(url: string) {
  let res = await instance.get(url, {
    responseType: "blob",
  });
  return res
}

export function postUpload(url: string, data = {}, onUploadProgress?: (progressEvent: any) => void) {
  return instance.post(url, data, {
    headers: {
      "Content-Type": "multipart/form-data",
      "X-Request-ID": `${generateRandomString(12)}`,
    },
    onUploadProgress,
  });
}

export function postChat(url: string, data = {}) {
  return instance.post(url, data, {
    headers: {
      "Content-Type": "text/event-stream;charset=utf-8",
      "X-Request-ID": `${generateRandomString(12)}`,
    },
  });
}

export function post(url: string, data = {}, config?: any) {
  return instance.post(url, data, config);
}

export function put(url: string, data = {}) {
  return instance.put(url, data);
}

export function del(url: string, data?: any) {
  return instance.delete(url, { data });
}
