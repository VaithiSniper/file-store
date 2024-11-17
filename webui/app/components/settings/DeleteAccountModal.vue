<script lang="ts" setup>
const model = defineModel({
  type: Boolean
})

const toast = useToast()

const loading = ref(false)

function onDelete() {
  loading.value = true

  setTimeout(() => {
    loading.value = false
    toast.add({icon: 'i-heroicons-check-circle', title: 'Your account has been deleted', color: 'red'})
    model.value = false
  }, 2000)
}
</script>

<template>
  <UDashboardModal
    v-model="model"
    :close-button="null"
    :ui="{
      icon: {
        base: 'text-red-500 dark:text-red-400'
      } as any,
      footer: {
        base: 'ml-16'
      } as any
    }"
    description="Are you sure you want to delete your account?"
    icon="i-heroicons-exclamation-circle"
    prevent-close
    title="Delete account"
  >
    <template #footer>
      <UButton
        :loading="loading"
        color="red"
        label="Delete"
        @click="onDelete"
      />
      <UButton
        color="white"
        label="Cancel"
        @click="model = false"
      />
    </template>
  </UDashboardModal>
</template>
