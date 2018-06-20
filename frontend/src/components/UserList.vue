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
    users: gql`
      {
        users
      }
    `,
  },
};
</script>
