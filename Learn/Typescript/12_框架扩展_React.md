# 第12章 框架扩展 - React + TypeScript

## 1. React + TypeScript 基础

### 1.1 项目搭建（Vite + React + TS）

```bash
# 使用 Vite 创建 React + TypeScript 项目
npm create vite@latest my-app -- --template react-ts

# 项目结构
# my-app/
# ├── src/
# │   ├── App.tsx
# │   ├── main.tsx
# │   └── vite-env.d.ts
# ├── tsconfig.json
# ├── tsconfig.node.json
# └── vite.config.ts
```

```jsonc
// tsconfig.json 关键配置
{
  "compilerOptions": {
    "target": "ES2020",
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "moduleResolution": "bundler",
    "jsx": "react-jsx",          // 使用新的 JSX 转换，无需 import React
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "skipLibCheck": true
  },
  "include": ["src"]
}
```

### 1.2 JSX 与 TSX

TSX 是 TypeScript 中的 JSX 语法扩展。与普通 TS 文件的关键区别：

```tsx
// 在 .tsx 文件中，尖括号类型断言会与 JSX 冲突
// 错误：编译器无法区分这是类型断言还是 JSX 元素
// const value = <string>someValue;

// 正确：在 .tsx 中使用 as 断言
const value = someValue as string;

// 泛型箭头函数需要加 extends 约束来消除歧义
// 错误：编译器会认为 <T> 是 JSX 标签
// const identity = <T>(arg: T): T => arg;

// 正确：添加约束
const identity = <T extends unknown>(arg: T): T => arg;
// 或者使用 trailing comma（某些配置下有效）
const identity2 = <T,>(arg: T): T => arg;
```
### 1.3 React.FC 类型（为什么现在不推荐用）

```tsx
// React.FC（FunctionComponent）的定义
type FC<P = {}> = FunctionComponent<P>;

interface FunctionComponent<P = {}> {
  (props: P): ReactNode;
  displayName?: string;
  defaultProps?: Partial<P>;  // 已废弃
}

// 使用 React.FC 的写法
const Button: React.FC<{ label: string }> = ({ label }) => {
  return <button>{label}</button>;
};

// 不推荐 React.FC 的原因：
// 1. 隐式包含 children（React 18 之前），容易导致类型不安全
// 2. 不支持泛型组件
// 3. defaultProps 处理有问题
// 4. 增加了不必要的类型复杂度

// 推荐写法：直接标注 props 类型
interface ButtonProps {
  label: string;
  variant?: 'primary' | 'secondary';
}

function Button({ label, variant = 'primary' }: ButtonProps) {
  return <button className={variant}>{label}</button>;
}

// 或者箭头函数
const Button = ({ label, variant = 'primary' }: ButtonProps) => {
  return <button className={variant}>{label}</button>;
};

// 如果需要显式标注返回类型
function Button({ label }: ButtonProps): React.ReactElement {
  return <button>{label}</button>;
}
```

---

## 2. 组件类型

### 2.1 函数组件的类型定义

```tsx
// 最简单的组件 - 无 props
function Welcome() {
  return <h1>Hello</h1>;
}

// 带 props 的组件
interface CardProps {
  title: string;
  children: React.ReactNode; // 最通用的 children 类型
  footer?: React.ReactElement; // 只接受 React 元素
}

function Card({ title, children, footer }: CardProps) {
  return (
    <div className="card">
      <h2>{title}</h2>
      <div>{children}</div>
      {footer && <div className="footer">{footer}</div>}
    </div>
  );
}
```

### 2.2 为什么不推荐 React.FC

```tsx
// React.FC 隐式包含 children，且泛型写法不够灵活
// 现代 React + TS 项目推荐直接标注 props 类型

// 不推荐
const Button: React.FC<ButtonProps> = ({ label }) => <button>{label}</button>;

// 推荐
function Button({ label }: ButtonProps) {
  return <button>{label}</button>;
}
```

## 3. Hooks 类型

```tsx
// useState - 大多数情况自动推断
const [count, setCount] = useState(0); // number
const [user, setUser] = useState<User | null>(null); // 需要显式标注

// useRef
const inputRef = useRef<HTMLInputElement>(null);
// inputRef.current 是 HTMLInputElement | null

// useReducer - 可辨识联合 Action
type Action =
  | { type: 'increment'; payload: number }
  | { type: 'decrement'; payload: number }
  | { type: 'reset' };

interface State { count: number }

function reducer(state: State, action: Action): State {
  switch (action.type) {
    case 'increment': return { count: state.count + action.payload };
    case 'decrement': return { count: state.count - action.payload };
    case 'reset': return { count: 0 };
  }
}

const [state, dispatch] = useReducer(reducer, { count: 0 });
dispatch({ type: 'increment', payload: 5 }); // 类型安全
```

## 4. 事件处理

```tsx
function Form() {
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    console.log(e.target.value);
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
  };

  const handleClick = (e: React.MouseEvent<HTMLButtonElement>) => {
    console.log(e.clientX, e.clientY);
  };

  return (
    <form onSubmit={handleSubmit}>
      <input onChange={handleChange} />
      <button onClick={handleClick}>Submit</button>
    </form>
  );
}
```

## 5. 高级模式

### 泛型组件

```tsx
// 通用列表组件
interface ListProps<T> {
  items: T[];
  renderItem: (item: T) => React.ReactNode;
  keyExtractor: (item: T) => string;
}

function List<T>({ items, renderItem, keyExtractor }: ListProps<T>) {
  return (
    <ul>
      {items.map(item => (
        <li key={keyExtractor(item)}>{renderItem(item)}</li>
      ))}
    </ul>
  );
}

// 使用时 T 自动推断
<List
  items={[{ id: '1', name: 'Alice' }]}
  renderItem={(user) => <span>{user.name}</span>}
  keyExtractor={(user) => user.id}
/>
```

### Polymorphic 组件（as prop）

```tsx
type PolymorphicProps<E extends React.ElementType> = {
  as?: E;
  children: React.ReactNode;
} & Omit<React.ComponentPropsWithoutRef<E>, 'as' | 'children'>;

function Box<E extends React.ElementType = 'div'>({
  as,
  children,
  ...props
}: PolymorphicProps<E>) {
  const Component = as || 'div';
  return <Component {...props}>{children}</Component>;
}

// 使用
<Box as="a" href="/about">Link</Box>  // 有 href 属性
<Box as="button" onClick={() => {}}>Button</Box> // 有 onClick
```

## 6. 状态管理

### Zustand + TypeScript

```typescript
import { create } from 'zustand';

interface TodoStore {
  todos: Todo[];
  addTodo: (text: string) => void;
  toggleTodo: (id: number) => void;
}

const useTodoStore = create<TodoStore>((set) => ({
  todos: [],
  addTodo: (text) => set((state) => ({
    todos: [...state.todos, { id: Date.now(), text, done: false }]
  })),
  toggleTodo: (id) => set((state) => ({
    todos: state.todos.map(t => t.id === id ? { ...t, done: !t.done } : t)
  })),
}));
```

---

**下一章：** [13_框架扩展_Node](./13_框架扩展_Node.md) - Node.js + TypeScript 后端开发
