const AuthPlugin = {
  install(Vue, options) {
    Vue.prototype.$setCurrentUser = function(user) {
      this.user = user;
    };
    Vue.prototype.$currentUser = function() {
      return this.user;
    };
  },
};

export { AuthPlugin };
