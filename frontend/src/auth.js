const AuthPlugin = {
  install(Vue, options) {
    Vue.prototype.$setCurrentUser = function (user) {
      Vue.prototype.user = user;
    };
    Vue.prototype.$currentUser = function () {
      return Vue.prototype.user;
    };
  },
};

export { AuthPlugin };
