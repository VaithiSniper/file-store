<script lang="ts" setup>
import {format, isToday} from 'date-fns'
import type {Mail} from '~/types'

defineProps({
  mail: {
    type: Object as PropType<Mail>,
    required: true
  },
  selected: {
    type: Boolean,
    default: false
  }
})
</script>

<template>
  <UDashboardPanelContent>
    <div class="flex justify-between">
      <div class="flex items-center gap-4">
        <UAvatar
          :alt="mail.from.name"
          size="lg"
          v-bind="mail.from.avatar"
        />

        <div class="min-w-0">
          <p class="text-gray-900 dark:text-white font-semibold">
            {{ mail.from.name }}
          </p>
          <p class="text-gray-500 dark:text-gray-400 font-medium">
            {{ mail.subject }}
          </p>
        </div>
      </div>

      <p class="font-medium text-gray-900 dark:text-white">
        {{
          isToday(new Date(mail.date)) ? format(new Date(mail.date), 'HH:mm') : format(new Date(mail.date), 'dd MMM')
        }}
      </p>
    </div>

    <UDivider class="my-5"/>

    <div class="flex-1">
      <p class="text-lg">
        {{ mail.body }}
      </p>
    </div>

    <UDivider class="my-5"/>

    <form @submit.prevent>
      <UTextarea
        :placeholder="`Reply to ${mail.from.name}`"
        :rows="5"
        color="gray"
        required
        size="xl"
      >
        <UButton
          class="absolute bottom-2.5 right-3.5"
          color="black"
          icon="i-heroicons-paper-airplane"
          label="Send"
          type="submit"
        />
      </UTextarea>
    </form>
  </UDashboardPanelContent>
</template>
