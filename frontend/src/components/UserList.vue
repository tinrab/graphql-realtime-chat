<template>
  <div>
    <h3>Users</h3>
    <div v-for="user of users"
         :key="user">
      {{user}}
    </div>
  </div>
</template>

<script>
import gql from 'graphql-tag';

export default {
  data() {
    return {
      users: [],
    };
  },
  apollo: {
    users() {
      const user = this.$currentUser();
      return {
        query: gql`
          {
            users
          }
        `,
        subscribeToMore: {
          document: gql`
            subscription($user: String!) {
              userJoined(user: $user)
            }
          `,
          variables: () => ({ user: user }),
          updateQuery: (prev, { subscriptionData }) => {
            if (!subscriptionData.data) {
              return prev;
            }
            const user = subscriptionData.data.userJoined;
            if (prev.users.find((u) => u === user)) {
              return prev;
            }
            return Object.assign({}, prev, {
              users: [user, ...prev.users],
            });
          },
        },
      };
    },
  },
};
</script>
