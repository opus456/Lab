# 13. 框架扩展：Node.js + TypeScript

## 1. Node.js + TypeScript 基础

### 1.1 项目搭建

```bash
# 初始化项目
mkdir my-node-ts-app && cd my-node-ts-app
npm init -y

# 安装 TypeScript 和 Node.js 类型
npm install -D typescript @types/node

# 初始化 tsconfig
npx tsc --init
```

### 1.2 @types/node

`@types/node` 提供了 Node.js 所有内置模块的类型定义：

```typescript
import fs from 'fs';
import path from 'path';
import { EventEmitter } from 'events';
import { IncomingMessage, ServerResponse } from 'http';

// fs 模块 - 类型自动推断
const content: string = fs.readFileSync('./file.txt', 'utf-8');

// path 模块
const fullPath: string = path.resolve(__dirname, 'src', 'index.ts');

// Buffer 类型
const buf: Buffer = Buffer.from('hello', 'utf-8');

// Stream 类型
import { Readable, Writable, Transform } from 'stream';

const readable: Readable = fs.createReadStream('./large-file.txt');
const writable: Writable = fs.createWriteStream('./output.txt');

// process 类型
const env: NodeJS.ProcessEnv = process.env;
const port: string | undefined = process.env.PORT;

// 事件类型
class MyEmitter extends EventEmitter {
  emit(event: 'data', payload: { id: number; name: string }): boolean;
  emit(event: string, ...args: any[]): boolean {
    return super.emit(event, ...args);
  }
}
```

### 1.3 tsconfig for Node.js

针对不同 Node.js 版本的推荐配置：

```jsonc
// tsconfig.json - Node.js 18+
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "lib": ["ES2022"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist"]
}
```

```jsonc
// tsconfig.json - Node.js 20+ (ESM 优先)
{
  "compilerOptions": {
    "target": "ES2023",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "lib": ["ES2023"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "exactOptionalPropertyTypes": true
  },
  "include": ["src/**/*"]
}
```

**关键配置说明：**

| 配置项 | 说明 |
|--------|------|
| `module: "NodeNext"` | 支持 ESM 和 CJS 混合使用 |
| `moduleResolution: "NodeNext"` | 遵循 Node.js 的模块解析规则 |
| `noUncheckedIndexedAccess` | 索引访问返回 `T \| undefined` |
| `exactOptionalPropertyTypes` | 区分 `undefined` 和可选属性 |

### 1.4 运行方式

#### 方式一：tsc + node（生产推荐）

```bash
# 编译后运行
npx tsc
node dist/index.js

# package.json scripts
{
  "scripts": {
    "build": "tsc",
    "start": "node dist/index.js",
    "dev": "tsc --watch & nodemon dist/index.js"
  }
}
```

#### 方式二：ts-node（开发调试）

```bash
npm install -D ts-node

# 直接运行 TypeScript
npx ts-node src/index.ts

# ESM 模式
npx ts-node --esm src/index.ts

# 使用 SWC 加速编译（推荐）
npm install -D @swc/core @swc/helpers
npx ts-node --swc src/index.ts
```

```jsonc
// tsconfig.json 中配置 ts-node
{
  "ts-node": {
    "swc": true,
    "transpileOnly": true,
    "esm": true
  }
}
```

#### 方式三：tsx（现代推荐）

```bash
npm install -D tsx

# 直接运行，无需配置
npx tsx src/index.ts

# 监听模式
npx tsx watch src/index.ts

# package.json
{
  "scripts": {
    "dev": "tsx watch src/index.ts",
    "start": "tsx src/index.ts"
  }
}
```

**三种方式对比：**

| 方式 | 速度 | 类型检查 | 适用场景 |
|------|------|----------|----------|
| tsc + node | 快（运行时） | 编译时检查 | 生产环境 |
| ts-node | 较慢 | 可选 | 开发调试 |
| tsx | 快 | 无 | 开发、脚本 |

---

## 2. Express + TypeScript

### 2.1 基础搭建与类型定义

```bash
npm install express
npm install -D @types/express
```

```typescript
import express, { Application, Request, Response, NextFunction } from 'express';

const app: Application = express();
app.use(express.json());

// 基础路由 - 类型自动推断
app.get('/health', (req: Request, res: Response) => {
  res.json({ status: 'ok', timestamp: Date.now() });
});

app.listen(3000, () => {
  console.log('Server running on port 3000');
});
```

### 2.2 Request/Response 类型扩展

```typescript
// 泛型参数定义请求的各个部分
// Request<Params, ResBody, ReqBody, Query>

interface CreateUserBody {
  name: string;
  email: string;
  age?: number;
}

interface UserResponse {
  id: string;
  name: string;
  email: string;
  createdAt: string;
}

interface UserParams {
  id: string;
}

interface UserQuery {
  page?: string;
  limit?: string;
  sort?: 'asc' | 'desc';
}

// 精确类型的路由处理器
app.post('/users',
  (req: Request<{}, UserResponse, CreateUserBody>, res: Response<UserResponse>) => {
    const { name, email, age } = req.body; // 类型安全
    const user: UserResponse = {
      id: crypto.randomUUID(),
      name,
      email,
      createdAt: new Date().toISOString(),
    };
    res.status(201).json(user);
  }
);
```

## 3. NestJS（TypeScript-first 框架）

NestJS 是专为 TypeScript 设计的 Node.js 框架，大量使用装饰器和依赖注入。

```typescript
// 控制器
import { Controller, Get, Post, Body, Param } from '@nestjs/common';

@Controller('users')
export class UserController {
  constructor(private readonly userService: UserService) {}

  @Get()
  findAll(): Promise<User[]> {
    return this.userService.findAll();
  }

  @Get(':id')
  findOne(@Param('id') id: string): Promise<User> {
    return this.userService.findOne(id);
  }

  @Post()
  create(@Body() createUserDto: CreateUserDto): Promise<User> {
    return this.userService.create(createUserDto);
  }
}

// DTO 与验证
import { IsString, IsEmail, IsOptional, MinLength } from 'class-validator';

class CreateUserDto {
  @IsString()
  @MinLength(2)
  name: string;

  @IsEmail()
  email: string;

  @IsOptional()
  @IsString()
  avatar?: string;
}
```

## 4. 数据库操作

### Prisma（推荐）

```typescript
// schema.prisma 定义模型后，Prisma 自动生成类型
const user = await prisma.user.findUnique({
  where: { id: '1' },
  include: { posts: true }, // 返回类型自动包含 posts
});
// user 的类型是 User & { posts: Post[] } | null

// 类型安全的查询
const users = await prisma.user.findMany({
  where: {
    email: { contains: '@example.com' },
    age: { gte: 18 },
  },
  select: { name: true, email: true }, // 返回类型只有 name 和 email
});
```

## 5. tRPC - 端到端类型安全

```typescript
// 服务端定义路由
import { initTRPC } from '@trpc/server';
import { z } from 'zod';

const t = initTRPC.create();

const appRouter = t.router({
  getUser: t.procedure
    .input(z.object({ id: z.string() }))
    .query(async ({ input }) => {
      return { id: input.id, name: 'Alice' };
    }),

  createUser: t.procedure
    .input(z.object({ name: z.string(), email: z.string().email() }))
    .mutation(async ({ input }) => {
      return { id: '1', ...input };
    }),
});

export type AppRouter = typeof appRouter;

// 客户端（自动获得完整类型推断，无需手动定义）
import { createTRPCClient } from '@trpc/client';
import type { AppRouter } from '../server';

const client = createTRPCClient<AppRouter>({ /* ... */ });
const user = await client.getUser.query({ id: '1' });
// user 自动推断为 { id: string; name: string }
```

## 6. 测试

```typescript
import { describe, it, expect, vi } from 'vitest';

// 类型安全的 Mock
interface UserRepository {
  findById(id: string): Promise<User | null>;
  save(user: User): Promise<User>;
}

const mockRepo: UserRepository = {
  findById: vi.fn().mockResolvedValue({ id: '1', name: 'Alice' }),
  save: vi.fn().mockImplementation(async (user) => user),
};

describe('UserService', () => {
  it('should find user by id', async () => {
    const service = new UserService(mockRepo);
    const user = await service.getUser('1');
    expect(user?.name).toBe('Alice');
  });
});
```

---

**下一章：** [14_框架扩展_Vue_Next](./14_框架扩展_Vue_Next.md) - Vue 3 和 Next.js 的 TypeScript 实践
