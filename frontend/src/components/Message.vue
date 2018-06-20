<template>
  <div>
    <strong>{{message.user}}</strong>: {{message.text}}
  </div>
</template>

<script>
export default {
  props: {
    message: {
      type: Object,
    },
  },
};
</script>
