import { createApp } from 'vue'
import App from './App.vue'
import './registerServiceWorker'
import {marked} from 'marked';

const markedMixin = {
    methods: {
         md: function (input: string) {
            
            return marked.parse(input);
        },
    },
};

createApp(App).mixin(markedMixin).mount('#app')
