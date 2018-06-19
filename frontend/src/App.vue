<template>
  <div>
    <div v-if="$apollo.loading">Loading...</div>

    <pre>{{JSON.stringify(users, 2, ' ')}}</pre>
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
    users: gql`
      {
        users
      }
    `,
  },
};
</script>
