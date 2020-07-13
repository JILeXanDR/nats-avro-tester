import Vue from 'vue';
import Vuetify from 'vuetify'
import App from './App.vue';
import 'vuetify/dist/vuetify.min.css'
import 'material-icons/css/material-icons.css';
import 'material-design-icons';
import Backend from "./backend";

Vue.use(Vuetify)

const vuetify = new Vuetify({});

Vue.prototype.$backend = new Backend(window.location.origin);

new Vue({
    vuetify,
    render: createElement => createElement(App),
}).$mount('#app');
