import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./router";
import "./assets/fonts.css";
import TDesign from "tdesign-vue-next";
// 引入组件库的少量全局样式变量
import "tdesign-vue-next/es/style/index.css";
import "@/assets/theme/theme.css";
import i18n from "./i18n";

// Fix passive event listener warnings
if (typeof window !== 'undefined') {
  const supportsPassive = (() => {
    let support = false;
    try {
      const opts = Object.defineProperty({}, 'passive', {
        get() {
          support = true;
        }
      });
      window.addEventListener('test', null as any, opts);
      window.removeEventListener('test', null as any, opts);
    } catch (e) {}
    return support;
  })();

  if (supportsPassive) {
    const addEvent = EventTarget.prototype.addEventListener;
    EventTarget.prototype.addEventListener = function(type: string, listener: any, options?: any) {
      const usesListenerOptions = typeof options === 'object' && options !== null;
      const useCapture = usesListenerOptions ? options.capture : options;
      
      options = usesListenerOptions ? options : {};
      options.passive = options.passive !== undefined ? options.passive : 
        ['touchstart', 'touchmove', 'wheel', 'mousewheel'].includes(type);
      options.capture = useCapture !== undefined ? useCapture : false;
      
      addEvent.call(this, type, listener, options);
    };
  }
}

const app = createApp(App);

app.use(TDesign);
app.use(createPinia());
app.use(router);
app.use(i18n);

app.mount("#app");
