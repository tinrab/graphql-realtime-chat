<template>
  <form class="col-12"
        v-on:submit.prevent="onPostClick">

    <div class="input-group">
      <input type="text"
             class="form-control"
             placeholder="Message..."
             v-model.trim="messageInput">
      <div class="input-group-append">
        <button class="btn btn-outline-secondary"
                type="submit">Post</button>
      </div>
    </div>

  </form>
</template>

<script>
import gql from 'graphql-tag';

export default {
  data() {
    return {
      messageInput: '',
    };
  },
  methods: {
    onPostClick() {
      const messageInput = this.messageInput;
      const user = this.$currentUser();

      this.$apollo
        .mutate({
          mutation: gql`
            mutation($user: String!, $text: String!) {
              postMessage(user: $user, text: $text) {
                id
              }
            }
          `,
          variables: {
            user: user,
            text: messageInput,
          },
        })
        .then((_) => {
          this.messageInput = '';
        })
        .catch((e) => {
          console.error(e);
        });
    },
  },
};
</script>
