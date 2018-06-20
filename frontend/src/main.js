import Vue from 'vue';
import { ApolloClient } from 'apollo-client';
import { HttpLink } from 'apollo-link-http';
import { InMemoryCache } from 'apollo-cache-inmemory';
import VueApollo from 'vue-apollo';
import { split } from 'apollo-link';
import { WebSocketLink } from 'apollo-link-ws';
import { getMainDefinition } from 'apollo-utilities';
import 'bootstrap/scss/bootstrap.scss';

import router from './router';
import App from './App.vue';
import { AuthPlugin } from './auth';

Vue.config.productionTip = false;

const httpLink = new HttpLink({
  uri: 'http://localhost:8080/graphql',
});
const wsLink = new WebSocketLink({
  uri: 'ws://localhost:8080/graphql',
  options: {
    reconnect: true,
  },
});
const link = split(
  ({ query }) => {
    const { kind, operation } = getMainDefinition(query);
    return kind === 'OperationDefinition' && operation === 'subscription';
  },
  wsLink,
  httpLink,
);
const apolloClient = new ApolloClient({
  link: link,
  cache: new InMemoryCache(),
});
const apolloProvider = new VueApollo({
  defaultClient: apolloClient,
});

Vue.use(VueApollo);
Vue.use(AuthPlugin);

const vm = new Vue({
  router,
  provide: apolloProvider.provide(),
  render: (h) => h(App),
});
vm.$mount('#app');
