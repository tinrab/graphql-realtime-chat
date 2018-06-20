<template>
  <div class="row">
    <form class="col-12 form"
          v-on:submit.prevent="onLoginClick">

      <div class="input-group">
        <input type="text"
               class="form-control"
               :class="{'is-invalid': !isInputValid}"
               placeholder="Username..."
               v-model.trim="userName"
               required>
        <div class="input-group-append">
          <button class="btn btn-secondary"
                  type="submit">Log in</button>
        </div>
      </div>

    </form>
  </div>
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
  methods: {
    onLoginClick() {
      if (this.isInputValid) {
        this.$setCurrentUser(this.userName);
        this.$router.push({ path: '/' });
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.form {
  max-width: 480px;
  margin: 0 auto;
}
</style>
