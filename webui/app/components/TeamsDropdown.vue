<script lang="ts" setup>
const teams = [{
  label: 'Nuxt',
  avatar: {
    src: 'https://avatars.githubusercontent.com/u/23360933?s=200&v=4'
  },
  click: () => {
    team.value = teams[0]
  }
}, {
  label: 'NuxtLabs',
  avatar: {
    src: 'https://avatars.githubusercontent.com/u/62017400?s=200&v=4'
  },
  click: () => {
    team.value = teams[1]
  }
}]

const actions = [{
  label: 'Create team',
  icon: 'i-heroicons-plus-circle'
}, {
  label: 'Manage teams',
  icon: 'i-heroicons-cog-8-tooth'
}]

const team = ref(teams[0])
</script>

<template>
  <UDropdown
    v-slot="{ open }"
    :items="[teams, actions]"
    :popper="{ strategy: 'absolute' }"
    :ui="{ width: 'w-full' }"
    class="w-full"
    mode="hover"
  >
    <UButton
      :class="[open && 'bg-gray-50 dark:bg-gray-800']"
      class="w-full"
      color="gray"
      variant="ghost"
    >
      <UAvatar
        :src="team.avatar.src"
        size="2xs"
      />

      <span class="truncate text-gray-900 dark:text-white font-semibold">{{ team.label }}</span>
    </UButton>
  </UDropdown>
</template>
