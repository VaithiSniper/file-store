<script lang="ts" setup>
import type {Mail} from '~/types'

const tabItems = [{
  label: 'All'
}, {
  label: 'Unread'
}]
const selectedTab = ref(0)

const dropdownItems = [[{
  label: 'Mark as unread',
  icon: 'i-heroicons-check-circle'
}, {
  label: 'Mark as important',
  icon: 'i-heroicons-exclamation-circle'
}], [{
  label: 'Star thread',
  icon: 'i-heroicons-star'
}, {
  label: 'Mute thread',
  icon: 'i-heroicons-pause-circle'
}]]

const {data: mails} = await useFetch<Mail[]>('/api/mails', {default: () => []})

// Filter mails based on the selected tab
const filteredMails = computed(() => {
  if (selectedTab.value === 1) {
    return mails.value.filter(mail => !!mail.unread)
  }

  return mails.value
})

const selectedMail = ref<Mail | null>()

const isMailPanelOpen = computed({
  get() {
    return !!selectedMail.value
  },
  set(value: boolean) {
    if (!value) {
      selectedMail.value = null
    }
  }
})

// Reset selected mail if it's not in the filtered mails
watch(filteredMails, () => {
  if (!filteredMails.value.find(mail => mail.id === selectedMail.value?.id)) {
    selectedMail.value = null
  }
})
</script>

<template>
  <UDashboardPage>
    <UDashboardPanel
      id="inbox"
      :resizable="{ min: 300, max: 500 }"
      :width="400"
    >
      <UDashboardNavbar
        :badge="filteredMails.length"
        title="Inbox"
      >
        <template #right>
          <UTabs
            v-model="selectedTab"
            :items="tabItems"
            :ui="{ wrapper: '', list: { height: 'h-9', tab: { height: 'h-7', size: 'text-[13px]' } } }"
          />
        </template>
      </UDashboardNavbar>

      <!-- ~/components/inbox/InboxList.vue -->
      <InboxList
        v-model="selectedMail"
        :mails="filteredMails"
      />
    </UDashboardPanel>

    <UDashboardPanel
      v-model="isMailPanelOpen"
      collapsible
      grow
      side="right"
    >
      <template v-if="selectedMail">
        <UDashboardNavbar>
          <template #toggle>
            <UDashboardNavbarToggle icon="i-heroicons-x-mark"/>

            <UDivider
              class="mx-1.5 lg:hidden"
              orientation="vertical"
            />
          </template>

          <template #left>
            <UTooltip text="Archive">
              <UButton
                color="gray"
                icon="i-heroicons-archive-box"
                variant="ghost"
              />
            </UTooltip>

            <UTooltip text="Move to junk">
              <UButton
                color="gray"
                icon="i-heroicons-archive-box-x-mark"
                variant="ghost"
              />
            </UTooltip>

            <UDivider
              class="mx-1.5"
              orientation="vertical"
            />

            <UPopover :popper="{ placement: 'bottom-start' }">
              <template #default="{ open }">
                <UTooltip
                  :prevent="open"
                  text="Snooze"
                >
                  <UButton
                    :class="[open && 'bg-gray-50 dark:bg-gray-800']"
                    color="gray"
                    icon="i-heroicons-clock"
                    variant="ghost"
                  />
                </UTooltip>
              </template>

              <template #panel="{ close }">
                <DatePicker @close="close"/>
              </template>
            </UPopover>
          </template>

          <template #right>
            <UTooltip text="Reply">
              <UButton
                color="gray"
                icon="i-heroicons-arrow-uturn-left"
                variant="ghost"
              />
            </UTooltip>

            <UTooltip text="Forward">
              <UButton
                color="gray"
                icon="i-heroicons-arrow-uturn-right"
                variant="ghost"
              />
            </UTooltip>

            <UDivider
              class="mx-1.5"
              orientation="vertical"
            />

            <UDropdown :items="dropdownItems">
              <UButton
                color="gray"
                icon="i-heroicons-ellipsis-vertical"
                variant="ghost"
              />
            </UDropdown>
          </template>
        </UDashboardNavbar>

        <!-- ~/components/inbox/InboxMail.vue -->
        <InboxMail :mail="selectedMail"/>
      </template>
      <div
        v-else
        class="flex-1 hidden lg:flex items-center justify-center"
      >
        <UIcon
          class="w-32 h-32 text-gray-400 dark:text-gray-500"
          name="i-heroicons-inbox"
        />
      </div>
    </UDashboardPanel>
  </UDashboardPage>
</template>
