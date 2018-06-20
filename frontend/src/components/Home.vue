<template>
  <div class="row">
    <app-user-list class="col-2" />
    <div class="col-10">

      <div class="row mb-3">
        <app-message-form />
      </div>

      <div class="row">
        <app-message-list class="col-12" />
      </div>

    </div>
  </div>
</template>

<script>
import UserList from '@/components/UserList.vue';
import MessageList from '@/components/MessageList.vue';
import MessageForm from '@/components/MessageForm.vue';

export default {
  components: {
    'app-user-list': UserList,
    'app-message-list': MessageList,
    'app-message-form': MessageForm,
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      if (!vm.$currentUser()) {
        vm.$router.push('login');
      }
    });
  },
};
</script>
