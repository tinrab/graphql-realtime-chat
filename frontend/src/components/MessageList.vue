<template>
  <div>
    <app-message v-for="message of messages"
                 :key="message.id"
                 :message="message">
    </app-message>
  </div>
</template>

<script>
import gql from 'graphql-tag';
import Message from '@/components/Message';

export default {
  components: {
    'app-message': Message,
  },
  data() {
    return {
      messages: [],
    };
  },
  apollo: {
    messages() {
      const user = this.$currentUser();
      return {
        query: gql`
          {
            messages {
              id
              user
              text
              createdAt
            }
          }
        `,
        subscribeToMore: {
          document: gql`
            subscription($user: String!) {
              messagePosted(user: $user) {
                id
                user
                text
                createdAt
              }
            }
          `,
          variables: () => ({ user: user }),
          updateQuery: (prev, { subscriptionData }) => {
            const message = subscriptionData.data.messagePosted;
            return Object.assign({}, prev, {
              messages: [message, ...prev.messages],
            });
          },
        },
      };
    },
  },
};
</script>
