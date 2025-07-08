import {
  create,
  NButton,
  NInput,
  NDatePicker,
  NForm,
  NFormItem,
  NTable,
  NCard,
  NSpace,
  NIcon,
  NConfigProvider,
  NDialogProvider,
  NMessageProvider,
  NModal,
  NSwitch,
  NDropdown
} from 'naive-ui'

export function createNaiveUI() {
  return create({
    components: [
      NButton,
      NInput,
      NDatePicker,
      NForm,
      NFormItem,
      NTable,
      NCard,
      NSpace,
      NIcon,
      NConfigProvider,
      NDialogProvider,
      NMessageProvider,
      NModal,
      NSwitch,
      NDropdown
    ]
  })
} 