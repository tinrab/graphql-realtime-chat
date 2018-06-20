<template>
  <form>

    <input type="text"
           class="form-control"
           :class="{'is-invalid': !isInputValid, 'mb-2': isInputValid}"
           placeholder="Username..."
           v-model.trim="userName"
           required>
    <div class="invalid-feedback mb-2">
      Please provide a valid username.
    </div>

  </form>
</template>

<script>
export default {
  data() {
    return {
      userName: '',
    };
  },
  computed: {
    isInputValid() {
      return this.userName.length > 0;
    },
  },
};
</script>

