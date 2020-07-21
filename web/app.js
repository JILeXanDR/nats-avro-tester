import Vue from 'vue';
import Vuetify from 'vuetify';
import VJsoneditor from 'v-jsoneditor/src/index';
import App from './App.vue';
import 'vuetify/dist/vuetify.min.css';
import Backend from "./backend";
import {NotifierPlugin} from "./plugins";
import VueMoment from 'vue-moment';

Vue.use(Vuetify);
Vue.use(NotifierPlugin);
Vue.use(VJsoneditor);
Vue.use(VueMoment);

const vuetify = new Vuetify({});

// TODO: replace with ENV vars in the future...
Vue.prototype.$backend = new Backend(window.location.origin);

new Vue({
    vuetify,
    render: createElement => createElement(App),
}).$mount('#app');
