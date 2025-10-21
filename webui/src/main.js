import { createApp, reactive } from "vue";
import App from "./App.vue";
import router from "./router";
import axios from "./services/axios.js";
import ToastPlugin from "vue-toast-notification";
import "vue-toast-notification/dist/theme-bootstrap.css";

// Icons
import { OhVueIcon, addIcons } from "oh-vue-icons";
import { BiInfoCircle } from "oh-vue-icons/icons";

addIcons(BiInfoCircle);

// Css
import "./assets/global.css";

const app = createApp(App);
app.config.globalProperties.$axios = axios;

// Global component registration (app.component("ComponentA", ComponentA))
app.component("v-icon", OhVueIcon);

app.use(router);
app.use(ToastPlugin);
app.mount("#app");
