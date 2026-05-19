# 第14章 框架扩展：Vue 3 + TypeScript & Next.js + TypeScript

## 前言

现代前端框架对 TypeScript 的支持已经从"可选"变为"推荐"甚至"默认"。Vue 3 从底层重写了类型系统，Next.js 则天然拥抱 TypeScript。本章将深入探讨如何在这两个主流框架中充分利用 TypeScript 的类型能力，构建类型安全的应用。

---

## 第一部分：Vue 3 + TypeScript

### 1. Vue 3 TypeScript 基础

#### 1.1 项目搭建（Vite + Vue + TS）

```bash
# 使用 Vite 创建 Vue + TypeScript 项目
npm create vite@latest my-vue-app -- --template vue-ts

# 或使用 create-vue（官方推荐）
npm create vue@latest
# 选择 TypeScript 支持
```

项目结构中关键的类型配置文件：

```json
// tsconfig.json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "ESNext",
    "moduleResolution": "bundler",
    "strict": true,
    "jsx": "preserve",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "esModuleInterop": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "skipLibCheck": true,
    "noEmit": true,
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["src/**/*.ts", "src/**/*.tsx", "src/**/*.vue"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

`env.d.ts` 文件让 TypeScript 识别 `.vue` 文件：

```typescript
/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}
```
#### 1.2 defineComponent 与 script setup

Vue 3 提供了两种 TypeScript 友好的组件编写方式：

**Options API + defineComponent：**

```typescript
import { defineComponent, PropType } from 'vue'

interface User {
  id: number
  name: string
  email: string
}

export default defineComponent({
  name: 'UserCard',
  props: {
    user: {
      type: Object as PropType<User>,
      required: true
    },
    showEmail: {
      type: Boolean,
      default: false
    }
  },
  data() {
    // this.user 被正确推断为 User 类型
    return {
      isEditing: false
    }
  },
  computed: {
    displayName(): string {
      return this.user.name.toUpperCase()
    }
  },
  methods: {
    toggleEdit(): void {
      this.isEditing = !this.isEditing
    }
  }
})
```

**Composition API + script setup（推荐）：**

```typescript
<script setup lang="ts">
import { ref, computed } from 'vue'

interface User {
  id: number
  name: string
  email: string
}

const props = defineProps<{
  user: User
  showEmail?: boolean
}>()

const isEditing = ref(false) // 自动推断为 Ref<boolean>

const displayName = computed(() => props.user.name.toUpperCase())
// displayName 类型为 ComputedRef<string>
</script>
```
#### 1.3 为什么 Vue 3 对 TypeScript 支持大幅改善

Vue 2 的类型支持存在根本性问题：

```typescript
// Vue 2 的问题：this 的类型推断困难
// Options API 中 this 指向组件实例，但 TypeScript 难以推断
// mixins 的类型几乎无法正确推断
// 模板中的类型检查完全缺失

// Vue 3 的改进：
// 1. 源码用 TypeScript 重写 —— 类型定义与源码同步
// 2. Composition API 天然适合类型推断 —— 函数返回值可推断
// 3. defineComponent 提供正确的 this 类型上下文
// 4. 泛型组件支持 —— <script setup generic="T">
// 5. Volar 提供模板内类型检查
```

**设计原因：** Composition API 的函数式设计让 TypeScript 的类型推断自然工作。不再需要通过 `this` 访问数据，避免了 Vue 2 中 `this` 类型推断的复杂性。

---

### 2. Composition API 类型

#### 2.1 ref() 和 reactive() 的类型推断

```typescript
import { ref, reactive, Ref } from 'vue'

// ref 自动推断类型
const count = ref(0)           // Ref<number>
const message = ref('hello')   // Ref<string>
const isActive = ref(true)     // Ref<boolean>

// 显式指定类型（当初始值不能表达完整类型时）
const user = ref<User | null>(null)  // Ref<User | null>
const items = ref<string[]>([])      // Ref<string[]>

// 复杂类型
interface FormState {
  username: string
  password: string
  remember: boolean
}

const form = ref<FormState>({
  username: '',
  password: '',
  remember: false
})

// reactive 自动推断
const state = reactive({
  count: 0,
  name: 'Vue'
})
// 类型为 { count: number; name: string }

// reactive 不推荐用泛型，因为返回类型会丢失响应式包装信息
// 推荐：直接传入对象让 TypeScript 推断
const complexState = reactive<{
  users: User[]
  loading: boolean
  error: string | null
}>({
  users: [],
  loading: false,
  error: null
})
```
**实际工程场景：** 表单状态管理

```typescript
<script setup lang="ts">
import { ref, reactive, watch } from 'vue'

interface ValidationRule {
  required?: boolean
  min?: number
  max?: number
  pattern?: RegExp
  message: string
}

interface FieldState<T> {
  value: T
  error: string | null
  touched: boolean
  rules: ValidationRule[]
}

function useField<T>(initialValue: T, rules: ValidationRule[] = []): FieldState<T> {
  return reactive({
    value: initialValue,
    error: null,
    touched: false,
    rules
  })
}

const email = useField('', [
  { required: true, message: '邮箱不能为空' },
  { pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/, message: '邮箱格式不正确' }
])

const age = useField<number | null>(null, [
  { required: true, message: '年龄不能为空' },
  { min: 0, max: 150, message: '年龄必须在0-150之间' }
])
</script>
```

#### 2.2 computed 类型

```typescript
import { ref, computed, ComputedRef } from 'vue'

const firstName = ref('张')
const lastName = ref('三')

// 自动推断返回类型为 ComputedRef<string>
const fullName = computed(() => `${firstName.value}${lastName.value}`)

// 可写计算属性
const fullNameWritable = computed({
  get(): string {
    return `${firstName.value}${lastName.value}`
  },
  set(newValue: string) {
    // 假设姓一个字，名可能多个字
    firstName.value = newValue.slice(0, 1)
    lastName.value = newValue.slice(1)
  }
})
```

### 3. 组件类型

```vue
<script setup lang="ts">
// defineProps - 泛型方式定义 Props
const props = defineProps<{
  title: string
  count?: number
  items: string[]
}>()

// 带默认值
const props2 = withDefaults(defineProps<{
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
}>(), {
  size: 'md',
  disabled: false,
})

// defineEmits - 类型安全的事件
const emit = defineEmits<{
  change: [value: string]
  submit: [data: FormData]
  'update:modelValue': [value: number]
}>()

emit('change', 'hello') // OK
emit('change', 123)     // Error!

// defineExpose - 暴露给父组件的方法
defineExpose({
  reset: () => { /* ... */ },
  validate: (): boolean => true,
})
</script>
```

### 4. Vue Router + TypeScript

```typescript
import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  { path: '/', component: () => import('./views/Home.vue') },
  { path: '/user/:id', component: () => import('./views/User.vue') },
]

// 在组件中使用
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const id = route.params.id // string | string[]

const router = useRouter()
router.push({ name: 'user', params: { id: '123' } })
```

### 5. Pinia + TypeScript

```typescript
import { defineStore } from 'pinia'

interface UserState {
  user: User | null
  token: string | null
  loading: boolean
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    user: null,
    token: null,
    loading: false,
  }),

  getters: {
    isLoggedIn: (state): boolean => !!state.token,
    displayName: (state): string => state.user?.name ?? 'Guest',
  },

  actions: {
    async login(email: string, password: string) {
      this.loading = true
      try {
        const { user, token } = await api.login(email, password)
        this.user = user
        this.token = token
      } finally {
        this.loading = false
      }
    },
  },
})
```

---

## 第二部分：Next.js + TypeScript

### 1. App Router 基础

```typescript
// app/page.tsx - 页面组件
export default function HomePage() {
  return <h1>Home</h1>
}

// app/layout.tsx - 布局
export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html><body>{children}</body></html>
  )
}

// Metadata 类型
import { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'My App',
  description: 'Built with Next.js',
}
```

### 2. 数据获取

```typescript
// Server Component 直接 async
interface Post { id: string; title: string; content: string }

export default async function PostPage({ params }: { params: { id: string } }) {
  const post = await fetch(`/api/posts/${params.id}`).then(r => r.json()) as Post
  return <article><h1>{post.title}</h1><p>{post.content}</p></article>
}

// Server Actions
'use server'

import { z } from 'zod'

const CreatePostSchema = z.object({
  title: z.string().min(1),
  content: z.string().min(10),
})

export async function createPost(formData: FormData) {
  const data = CreatePostSchema.parse({
    title: formData.get('title'),
    content: formData.get('content'),
  })
  await db.post.create({ data })
}
```

### 3. API Routes

```typescript
// app/api/users/route.ts
import { NextRequest, NextResponse } from 'next/server'

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams
  const page = Number(searchParams.get('page') ?? '1')

  const users = await db.user.findMany({ skip: (page - 1) * 10, take: 10 })
  return NextResponse.json(users)
}

export async function POST(request: NextRequest) {
  const body = await request.json()
  const user = await db.user.create({ data: body })
  return NextResponse.json(user, { status: 201 })
}
```

### 4. 全栈类型安全

```typescript
// 使用 tRPC + Next.js 实现端到端类型安全
// 服务端定义 → 客户端自动获得类型
// 修改 API 返回值 → 前端立即得到类型错误
// 零手动类型定义，零运行时开销
```

---

**总结：** 本系列文档从 TypeScript 基础到框架实战，覆盖了完整的学习路径。建议按照 README 中的学习路线循序渐进，每个章节的代码示例都可以直接运行验证。
