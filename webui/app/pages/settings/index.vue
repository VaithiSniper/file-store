<script lang="ts" setup>
import type {FormError, FormSubmitEvent} from '#ui/types'

const fileRef = ref<HTMLInputElement>()
const isDeleteAccountModalOpen = ref(false)

const state = reactive({
  name: 'Benjamin Canac',
  email: 'ben@nuxtlabs.com',
  username: 'benjamincanac',
  avatar: '',
  bio: '',
  password_current: '',
  password_new: ''
})

const toast = useToast()

function validate(state: any): FormError[] {
  const errors = []
  if (!state.name) errors.push({path: 'name', message: 'Please enter your name.'})
  if (!state.email) errors.push({path: 'email', message: 'Please enter your email.'})
  if ((state.password_current && !state.password_new) || (!state.password_current && state.password_new)) errors.push({
    path: 'password',
    message: 'Please enter a valid password.'
  })
  return errors
}

function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement

  if (!input.files?.length) {
    return
  }

  state.avatar = URL.createObjectURL(input.files[0])
}

function onFileClick() {
  fileRef.value?.click()
}

async function onSubmit(event: FormSubmitEvent<any>) {
  // Do something with data
  console.log(event.data)

  toast.add({title: 'Profile updated', icon: 'i-heroicons-check-circle'})
}
</script>

<template>
  <UDashboardPanelContent class="pb-24">
    <UDashboardSection
      description="Customize the look and feel of your dashboard."
      title="Theme"
    >
      <template #links>
        <UColorModeSelect color="gray"/>
      </template>
    </UDashboardSection>

    <UDivider class="mb-4"/>

    <UForm
      :state="state"
      :validate="validate"
      :validate-on="['submit']"
      @submit="onSubmit"
    >
      <UDashboardSection
        description="This information will be displayed publicly so be careful what you share."
        title="Profile"
      >
        <template #links>
          <UButton
            color="black"
            label="Save changes"
            type="submit"
          />
        </template>

        <UFormGroup
          :ui="{ container: '' }"
          class="grid grid-cols-2 gap-2 items-center"
          description="Will appear on receipts, invoices, and other communication."
          label="Name"
          name="name"
          required
        >
          <UInput
            v-model="state.name"
            autocomplete="off"
            icon="i-heroicons-user"
            size="md"
          />
        </UFormGroup>

        <UFormGroup
          :ui="{ container: '' }"
          class="grid grid-cols-2 gap-2"
          description="Used to sign in, for email receipts and product updates."
          label="Email"
          name="email"
          required
        >
          <UInput
            v-model="state.email"
            autocomplete="off"
            icon="i-heroicons-envelope"
            size="md"
            type="email"
          />
        </UFormGroup>

        <UFormGroup
          :ui="{ container: '' }"
          class="grid grid-cols-2 gap-2"
          description="Your unique username for logging in and your profile URL."
          label="Username"
          name="username"
          required
        >
          <UInput
            v-model="state.username"
            autocomplete="off"
            input-class="ps-[77px]"
            size="md"
            type="username"
          >
            <template #leading>
              <span class="text-gray-500 dark:text-gray-400 text-sm">nuxt.com/</span>
            </template>
          </UInput>
        </UFormGroup>

        <UFormGroup
          :ui="{ container: 'flex flex-wrap items-center gap-3', help: 'mt-0' }"
          class="grid grid-cols-2 gap-2"
          help="JPG, GIF or PNG. 1MB Max."
          label="Avatar"
          name="avatar"
        >
          <UAvatar
            :alt="state.name"
            :src="state.avatar"
            size="lg"
          />

          <UButton
            color="white"
            label="Choose"
            size="md"
            @click="onFileClick"
          />

          <input
            ref="fileRef"
            accept=".jpg, .jpeg, .png, .gif"
            class="hidden"
            type="file"
            @change="onFileChange"
          >
        </UFormGroup>

        <UFormGroup
          :ui="{ container: '' }"
          class="grid grid-cols-2 gap-2"
          description="Brief description for your profile. URLs are hyperlinked."
          label="Bio"
          name="bio"
        >
          <UTextarea
            v-model="state.bio"
            :rows="5"
            autoresize
            size="md"
          />
        </UFormGroup>

        <UFormGroup
          :ui="{ container: '' }"
          class="grid grid-cols-2 gap-2"
          description="Confirm your current password before setting a new one."
          label="Password"
          name="password"
        >
          <UInput
            id="password"
            v-model="state.password_current"
            placeholder="Current password"
            size="md"
            type="password"
          />
          <UInput
            id="password_new"
            v-model="state.password_new"
            class="mt-2"
            placeholder="New password"
            size="md"
            type="password"
          />
        </UFormGroup>
      </UDashboardSection>
    </UForm>

    <UDivider class="mb-4"/>

    <UDashboardSection
      description="No longer want to use our service? You can delete your account here. This action is not reversible. All information related to this account will be deleted permanently."
      title="Account"
    >
      <div>
        <UButton
          color="red"
          label="Delete account"
          size="md"
          @click="isDeleteAccountModalOpen = true"
        />
      </div>
    </UDashboardSection>

    <!-- ~/components/settings/DeleteAccountModal.vue -->
    <SettingsDeleteAccountModal v-model="isDeleteAccountModalOpen"/>
  </UDashboardPanelContent>
</template>
