export function NotifierPlugin(Vue, options) {
    Vue.mixin({
        data() {
            return {
                notifierPlugin: {
                    model: false,
                    timeout: 5000,
                    color: 'info',
                    text: '',
                },
            };
        },
    });

    const showNotification = function (type, message) {
        this.notifierPlugin.model = true;
        this.notifierPlugin.color = type;
        this.notifierPlugin.text = message;
    };

    // decorator
    const withTimeout = (fn, duration) => {
        return () => {
            setTimeout(() => fn(), duration);
        };
    };

    Vue.prototype.$notify = function (type, message) {
        let fn = () => showNotification.apply(this, [type, message]);
        if (this.notifierPlugin.model) {
            this.notifierPlugin.model = false;
            fn = withTimeout(fn, 100)
        }
        fn();
    }
}
