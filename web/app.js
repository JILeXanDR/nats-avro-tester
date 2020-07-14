import Vue from 'vue';
import Vuetify from 'vuetify'
import App from './App.vue';
import 'vuetify/dist/vuetify.min.css'
import Backend from "./backend";

Vue.use(Vuetify)

const vuetify = new Vuetify({});

// TODO: replace with ENV vars in the future...
Vue.prototype.$backend = new Backend(window.location.origin);

new Vue({
    vuetify,
    render: createElement => createElement(App),
}).$mount('#app');
