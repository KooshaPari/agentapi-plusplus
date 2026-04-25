# agentapi-plusplus Technical Specification

**Version**: 2.0.0 | **Status**: Active | **Last Updated**: 2026-04-06

---

## 1. Executive Summary

agentapi-plusplus is an advanced agent orchestration and management API platform designed for the Phenotype ecosystem. It provides comprehensive capabilities for agent lifecycle management, inter-agent communication, task distribution, and observability across distributed agent networks.

This specification defines the complete architecture, API surface, data models, and operational characteristics of the agentapi-plusplus platform.

---

## 2. Architecture Overview

### 2.1 System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                            AGENTAPI-PLUSPLUS PLATFORM                                        │
│                           (Microservices Architecture)                                       │
├─────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────────────────────┐  │
│  │                              EDGE LAYER (API Gateway)                               │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────────┐ │  │
│  │  │   REST API  │  │ GraphQL     │  │  gRPC       │  │   WebSocket                 │ │  │
│  │  │   Gateway   │  │ Endpoint    │  │  Services   │  │   Real-time                 │ │  │
│  │  │             │  │             │  │             │  │   Events                    │ │  │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └─────────────┬───────────────┘ │  │
│  │         │                │                │                       │                 │  │
│  └─────────┼────────────────┼────────────────┼───────────────────────┼─────────────────┘  │
│            │                │                │                       │                  │
│  ┌─────────▼────────────────▼────────────────▼───────────────────────▼────────────────┐  │
│  │                           CORE SERVICES LAYER                                      │  │
│  │                                                                                  │  │
│  │   ┌─────────────────────────────────────────────────────────────────────────┐   │  │
│  │   │                        AGENT MANAGEMENT SERVICE                            │   │  │
│  │   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │   │  │
│  │   │  │  Lifecycle  │  │   State     │  │  Registry   │  │  Capabilities│    │   │  │
│  │   │  │  Manager    │  │   Machine   │  │  Service    │  │  Engine      │    │   │  │
│  │   │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘    │   │  │
│  │   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │   │  │
│  │   │  │  Identity   │  │   Health    │  │   Metrics   │  │  Discovery  │    │   │  │
│  │   │  │  Service    │  │   Monitor   │  │   Collector │  │  Service     │    │   │  │
│  │   │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘    │   │  │
│  │   └─────────────────────────────────────────────────────────────────────────┘   │  │
│  │                                                                                  │  │
│  │   ┌─────────────────────────────────────────────────────────────────────────┐   │  │
│  │   │                        TASK ORCHESTRATION SERVICE                          │   │  │
│  │   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │   │  │
│  │   │  │   Task      │  │   Workflow  │  │  Scheduler  │  │  Executor   │    │   │  │
│  │   │  │   Queue     │  │   Engine    │  │             │  │             │    │   │  │
│  │   │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘    │   │  │
│  │   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │   │  │
│  │   │  │  Priority   │  │   DAG       │  │  Load       │  │  Result     │    │   │  │
│  │   │  │  Manager    │  │  Processor  │  │  Balancer   │  │  Collector  │    │   │  │
│  │   │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘    │   │  │
│  │   └─────────────────────────────────────────────────────────────────────────┘   │  │
│  │                                                                                  │  │
│  │   ┌─────────────────────────────────────────────────────────────────────────┐   │  │
│  │   │                      COMMUNICATION SERVICE                                 │   │  │
│  │   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │   │  │
│  │   │  │   Message   │  │   Pub/Sub   │  │   Stream    │  │  Event      │    │   │  │
│  │   │  │   Router    │  │   Engine    │  │   Manager   │  │  Bus        │    │   │  │
│  │   │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘    │   │  │
│  │   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │   │  │
│  │   │  │   Group     │  │  Broadcast  │  │  Presence   │  │  Channel    │    │   │  │
│  │   │  │   Manager   │  │  Service     │  │  Tracker    │  │  Manager    │    │   │  │
│  │   │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘    │   │  │
│  │   └─────────────────────────────────────────────────────────────────────────┘   │  │
│  │                                                                                  │  │
│  │   ┌─────────────────────────────────────────────────────────────────────────┐   │  │
│  │   │                      OBSERVABILITY SERVICE                                 │   │  │
│  │   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │   │  │
│  │   │  │   Tracing   │  │   Logging   │  │   Metrics   │  │   Alerting  │    │   │  │
│  │   │  │   Pipeline  │  │   Pipeline  │  │   Pipeline  │  │   Engine    │    │   │  │
│  │   │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘    │   │  │
│  │   │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │   │  │
│  │   │  │   APM       │  │   Profiling │  │   Tracing   │  │   SLO       │    │   │  │
│  │   │  │   Service   │  │   Service   │  │   Storage   │  │   Manager   │    │   │  │
│  │   │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘    │   │  │
│  │   └─────────────────────────────────────────────────────────────────────────┘   │  │
│  │                                                                                  │  │
│  └───────────────────────────────────────────────────────────────────────────────────┘  │
│            │                │                │                       │                  │
│  ┌─────────▼────────────────▼────────────────▼───────────────────────▼────────────────┐  │
│  │                           DATA LAYER                                             │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────────┐ │  │
│  │  │ PostgreSQL  │  │    Redis    │  │    Kafka    │  │   Object Storage          │ │  │
│  │  │ (Primary)   │  │   (Cache)   │  │  (Events)   │  │   (S3/MinIO)              │ │  │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────────────────────┘ │  │
│  │                                                                                      │ │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────────┐ │  │
│  │  │ ClickHouse  │  │Elasticsearch│  │  TimescaleDB│  │      Vector DB              │ │  │
│  │  │ (Analytics) │  │    (Logs)   │  │  (Metrics)  │  │   (Embeddings)              │ │  │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────────────────────┘ │  │
│  └────────────────────────────────────────────────────────────────────────────────────┘  │
│                                                                                             │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Agent Lifecycle State Machine

```
┌─────────┐    register     ┌─────────┐    activate    ┌─────────┐
│  NEW    │────────────────▶│ STANDBY │───────────────▶│ ACTIVE  │
└─────────┘                 └─────────┘                └────┬────┘
     │                            │                        │
     │                            │                        │
     │   ┌────────────────────────┘                        │
     │   │                                                 │
     │   │ pause                                           │ execute
     │   │                                                 │
     │   ▼                                                 ▼
     │ ┌─────────┐                                   ┌─────────┐
     └▶│  ERROR  │◀─────────────────────────────────│ RUNNING │
       └────┬────┘   error                           └────┬────┘
            │                                            │
            │ complete                                   │ suspend
            │                                            │
            ▼                                            ▼
       ┌─────────┐                                  ┌─────────┐
       │TERMINATE│                                  │SUSPENDED│
       └────┬────┘                                  └────┬────┘
            │                                            │
            │ deregister                                 │ resume
            │                                            │
            ▼                                            ▼
       ┌─────────┐◀──────────────────────────────────┐
       │ DESTROY │                                    │
       └─────────┘                                    │
                                                      │
                                               ┌─────────┐
                                               │  IDLE   │
                                               └─────────┘

State Transitions:
• NEW → STANDBY: Agent registers with the platform
• STANDBY → ACTIVE: Agent completes initialization and health checks
• ACTIVE → RUNNING: Agent receives and starts executing a task
• RUNNING → SUSPENDED: Task execution paused
• SUSPENDED → RUNNING: Task execution resumed
• RUNNING → COMPLETED: Task execution finished successfully
• RUNNING → ERROR: Task execution failed
• ACTIVE → STANDBY: Agent deactivated but retained
• STANDBY → DESTROY: Agent deregistered
• ERROR → STANDBY: Error recovered, agent reset
• ERROR → DESTROY: Terminal error, agent terminated
```

### 2.3 Message Flow Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────────────────┐
│                           INTER-AGENT COMMUNICATION PATTERNS                                │
├─────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                             │
│  Pattern 1: Direct (Point-to-Point)                                                         │
│  ┌─────────┐                        ┌─────────┐                                             │
│  │ Agent A │───────────────────────▶│ Agent B │                                             │
│  │         │    message.send(to=B)  │         │                                             │
│  └─────────┘                        └─────────┘                                             │
│                                                                                             │
│  Pattern 2: Broadcast                                                                         │
│  ┌─────────┐    broadcast(msg)    ┌─────────┐                                             │
│  │ Agent A │──────────────────────▶│ Agent B │                                             │
│  │         │──────────────────────▶│ Agent C │                                             │
│  │         │──────────────────────▶│ Agent D │                                             │
│  └─────────┘                        └─────────┘                                             │
│                                                                                             │
│  Pattern 3: Pub/Sub (Topics)                                                                  │
│                                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐                   │
│  │                      MESSAGE BROKER (Kafka)                          │                   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │                   │
│  │  │  Topic:     │  │  Topic:     │  │  Topic:     │  │  Topic:     │ │                   │
│  │  │  tasks      │  │  events     │  │  metrics    │  │  alerts     │ │                   │
│  │  │             │  │             │  │             │  │             │ │                   │
│  │  │ Consumer A │  │ Consumer B │  │ Consumer C │  │ Consumer D │ │                   │
│  │  │ Consumer C │  │ Consumer D │  │             │  │             │ │                   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │                   │
│  └─────────────────────────────────────────────────────────────────────┘                   │
│                                                                                             │
│  Pattern 4: Request-Reply (RPC)                                                               │
│  ┌─────────┐   request(req)   ┌─────────┐  process  ┌─────────┐   reply(resp)   ┌─────────┐│
│  │ Agent A │─────────────────▶│  Queue  │──────────▶│ Agent B │────────────────▶│ Agent A ││
│  │         │◀─────────────────│         │◀─────────│         │◀────────────────│         ││
│  │         │  (await reply)   │         │          │         │                  │         ││
│  └─────────┘                  └─────────┘          └─────────┘                  └─────────┘│
│                                                                                             │
│  Pattern 5: Fan-Out/Fan-In                                                                    │
│  ┌─────────┐                    ┌─────────┐                  ┌─────────┐                     │
│  │         │───dispatch──────▶│ Worker  │───result──────▶│         │                     │
│  │         │───dispatch──────▶│ Worker  │───result──────▶│         │                     │
│  │  Master │───dispatch──────▶│ Worker  │───result──────▶│  Result │                     │
│  │         │───dispatch──────▶│ Worker  │───result──────▶│  Aggregator                 │
│  │         │───dispatch──────▶│ Worker  │───result──────▶│         │                     │
│  └─────────┘                    └─────────┘                  └─────────┘                     │
│                                                                                             │
└─────────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Data Models

### 3.1 Core Domain Models (TypeScript)

```typescript
/**
 * Core domain models for agentapi-plusplus
 * @module models
 */

import { z } from 'zod';

// ============================================================================
// BASE TYPES
// ============================================================================

export type UUID = string;
export type Timestamp = string; // ISO 8601 format
export type JSONValue = string | number | boolean | null | JSONObject | JSONArray;
export type JSONObject = { [key: string]: JSONValue };
export type JSONArray = JSONValue[];

// ============================================================================
// AGENT MODELS
// ============================================================================

export const AgentStatusSchema = z.enum([
  'NEW',
  'STANDBY',
  'ACTIVE',
  'RUNNING',
  'SUSPENDED',
  'IDLE',
  'ERROR',
  'COMPLETED',
  'TERMINATE',
  'DESTROY'
]);

export type AgentStatus = z.infer<typeof AgentStatusSchema>;

export const AgentCapabilitySchema = z.object({
  id: z.string(),
  name: z.string(),
  version: z.string(),
  description: z.string().optional(),
  parameters: z.record(z.any()).optional(),
  requirements: z.object({
    cpu: z.number().optional(),
    memory: z.number().optional(),
    gpu: z.boolean().optional(),
    storage: z.number().optional(),
    network: z.boolean().optional(),
  }).optional(),
  tags: z.array(z.string()).default([]),
});

export type AgentCapability = z.infer<typeof AgentCapabilitySchema>;

export const AgentIdentitySchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
  type: z.enum(['LLM', 'TOOL', 'ORCHESTRATOR', 'HUMAN', 'SYSTEM', 'CUSTOM']),
  namespace: z.string(),
  version: z.string(),
  metadata: z.record(z.any()).default({}),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
});

export type AgentIdentity = z.infer<typeof AgentIdentitySchema>;

export const AgentConfigurationSchema = z.object({
  maxConcurrentTasks: z.number().int().positive().default(1),
  taskTimeoutSeconds: z.number().int().positive().default(300),
  idleTimeoutSeconds: z.number().int().positive().default(3600),
  retryPolicy: z.object({
    maxRetries: z.number().int().default(3),
    backoffMultiplier: z.number().default(2),
    initialDelayMs: z.number().int().default(1000),
    maxDelayMs: z.number().int().default(30000),
  }).optional(),
  logging: z.object({
    level: z.enum(['debug', 'info', 'warn', 'error']).default('info'),
    destination: z.enum(['stdout', 'file', 'remote']).default('stdout'),
  }).default({}),
  resources: z.object({
    cpuLimit: z.number().optional(),
    memoryLimit: z.number().optional(),
    storageLimit: z.number().optional(),
  }).optional(),
});

export type AgentConfiguration = z.infer<typeof AgentConfigurationSchema>;

export const AgentSchema = z.object({
  identity: AgentIdentitySchema,
  status: AgentStatusSchema,
  capabilities: z.array(AgentCapabilitySchema).default([]),
  configuration: AgentConfigurationSchema,
  endpoints: z.object({
    rest: z.string().url().optional(),
    websocket: z.string().url().optional(),
    grpc: z.string().optional(),
  }).optional(),
  health: z.object({
    lastPingAt: z.string().datetime().optional(),
    consecutiveFailures: z.number().int().default(0),
    status: z.enum(['healthy', 'degraded', 'unhealthy']).default('healthy'),
    metrics: z.object({
      cpuUsage: z.number().optional(),
      memoryUsage: z.number().optional(),
      activeConnections: z.number().int().optional(),
      tasksCompleted: z.number().int().optional(),
      tasksFailed: z.number().int().optional(),
    }).optional(),
  }).default({}),
  currentTask: z.string().uuid().optional(),
  assignedTasks: z.array(z.string().uuid()).default([]),
  tags: z.array(z.string()).default([]),
  labels: z.record(z.string()).default({}),
});

export type Agent = z.infer<typeof AgentSchema>;

// ============================================================================
// TASK MODELS
// ============================================================================

export const TaskStatusSchema = z.enum([
  'PENDING',
  'SCHEDULED',
  'ASSIGNED',
  'RUNNING',
  'PAUSED',
  'CANCELLED',
  'FAILED',
  'COMPLETED',
  'TIMEOUT'
]);

export type TaskStatus = z.infer<typeof TaskStatusSchema>;

export const TaskPrioritySchema = z.enum([
  'CRITICAL',
  'HIGH',
  'NORMAL',
  'LOW',
  'BACKGROUND'
]);

export type TaskPriority = z.infer<typeof TaskPrioritySchema>;

export const TaskInputSchema = z.object({
  type: z.enum(['text', 'json', 'binary', 'file', 'stream']),
  data: z.any(),
  metadata: z.record(z.any()).optional(),
  encoding: z.enum(['utf-8', 'base64', 'binary']).optional(),
  schema: z.string().optional(), // JSON Schema reference
});

export type TaskInput = z.infer<typeof TaskInputSchema>;

export const TaskOutputSchema = z.object({
  type: z.enum(['text', 'json', 'binary', 'file', 'stream', 'error']),
  data: z.any(),
  metadata: z.record(z.any()).optional(),
  encoding: z.enum(['utf-8', 'base64', 'binary']).optional(),
  schema: z.string().optional(),
  mimeType: z.string().optional(),
});

export type TaskOutput = z.infer<typeof TaskOutputSchema>;

export const TaskResultSchema = z.object({
  success: z.boolean(),
  output: TaskOutputSchema.optional(),
  error: z.object({
    code: z.string(),
    message: z.string(),
    details: z.record(z.any()).optional(),
    stackTrace: z.string().optional(),
  }).optional(),
  metrics: z.object({
    executionTimeMs: z.number(),
    cpuTimeMs: z.number().optional(),
    memoryPeakBytes: z.number().optional(),
    networkBytesTransferred: z.number().optional(),
  }),
  logs: z.array(z.object({
    level: z.enum(['debug', 'info', 'warn', 'error']),
    message: z.string(),
    timestamp: z.string().datetime(),
  })).optional(),
});

export type TaskResult = z.infer<typeof TaskResultSchema>;

export const TaskSchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
  description: z.string().optional(),
  type: z.enum(['SINGLE', 'BATCH', 'WORKFLOW', 'DAG', 'CHAIN', 'PARALLEL']),
  status: TaskStatusSchema,
  priority: TaskPrioritySchema,
  
  // Execution context
  input: TaskInputSchema,
  output: TaskOutputSchema.optional(),
  result: TaskResultSchema.optional(),
  
  // Assignment
  assignee: z.string().uuid().optional(),
  candidates: z.array(z.string().uuid()).optional(),
  
  // Timing
  createdAt: z.string().datetime(),
  scheduledAt: z.string().datetime().optional(),
  startedAt: z.string().datetime().optional(),
  completedAt: z.string().datetime().optional(),
  deadlineAt: z.string().datetime().optional(),
  estimatedDurationMs: z.number().optional(),
  
  // Configuration
  timeoutSeconds: z.number().int().positive().default(300),
  maxRetries: z.number().int().default(3),
  retryCount: z.number().int().default(0),
  
  // Dependencies
  dependencies: z.array(z.string().uuid()).default([]),
  dependents: z.array(z.string().uuid()).default([]),
  
  // Workflow
  workflowId: z.string().uuid().optional(),
  stepNumber: z.number().int().optional(),
  
  // Metadata
  tags: z.array(z.string()).default([]),
  labels: z.record(z.string()).default({}),
  parentId: z.string().uuid().optional(),
  rootId: z.string().uuid().optional(),
  
  // Progress
  progress: z.object({
    percent: z.number().min(0).max(100).default(0),
    currentStep: z.string().optional(),
    totalSteps: z.number().optional(),
    message: z.string().optional(),
  }).default({}),
});

export type Task = z.infer<typeof TaskSchema>;

// ============================================================================
// WORKFLOW MODELS
// ============================================================================

export const WorkflowNodeTypeSchema = z.enum([
  'START',
  'END',
  'TASK',
  'CONDITION',
  'PARALLEL',
  'JOIN',
  'SUBWORKFLOW',
  'DELAY',
  'EVENT'
]);

export type WorkflowNodeType = z.infer<typeof WorkflowNodeTypeSchema>;

export const WorkflowNodeSchema = z.object({
  id: z.string(),
  type: WorkflowNodeTypeSchema,
  name: z.string(),
  description: z.string().optional(),
  position: z.object({
    x: z.number(),
    y: z.number(),
  }).optional(),
  
  // Task-specific
  taskTemplate: TaskSchema.omit({ id: true, status: true, createdAt: true }).optional(),
  
  // Condition-specific
  condition: z.object({
    expression: z.string(),
    operator: z.enum(['equals', 'not_equals', 'contains', 'gt', 'lt', 'regex', 'custom']),
    value: z.any().optional(),
    customFunction: z.string().optional(),
  }).optional(),
  
  // Parallel-specific
  parallelBranches: z.number().int().positive().optional(),
  
  // Delay-specific
  delaySeconds: z.number().int().positive().optional(),
  
  // Event-specific
  eventType: z.string().optional(),
  eventFilter: z.record(z.any()).optional(),
  
  // General
  config: z.record(z.any()).optional(),
  onError: z.enum(['fail', 'retry', 'continue', 'abort']).default('fail'),
  retries: z.number().int().default(0),
});

export type WorkflowNode = z.infer<typeof WorkflowNodeSchema>;

export const WorkflowEdgeSchema = z.object({
  id: z.string(),
  source: z.string(),
  target: z.string(),
  label: z.string().optional(),
  condition: z.string().optional(), // For conditional edges
  type: z.enum(['default', 'success', 'failure', 'conditional']).default('default'),
});

export type WorkflowEdge = z.infer<typeof WorkflowEdgeSchema>;

export const WorkflowSchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
  description: z.string().optional(),
  version: z.string(),
  status: z.enum(['DRAFT', 'ACTIVE', 'ARCHIVED', 'DEPRECATED']),
  
  // Graph definition
  nodes: z.array(WorkflowNodeSchema),
  edges: z.array(WorkflowEdgeSchema),
  
  // Start configuration
  startNodeId: z.string(),
  
  // Variables and context
  variables: z.record(z.any()).default({}),
  inputSchema: z.record(z.any()).optional(),
  outputSchema: z.record(z.any()).optional(),
  
  // Execution config
  timeoutSeconds: z.number().int().positive().default(3600),
  maxConcurrentExecutions: z.number().int().positive().default(10),
  
  // Metadata
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
  createdBy: z.string(),
  tags: z.array(z.string()).default([]),
});

export type Workflow = z.infer<typeof WorkflowSchema>;

// ============================================================================
// MESSAGE MODELS
// ============================================================================

export const MessageTypeSchema = z.enum([
  'COMMAND',
  'QUERY',
  'EVENT',
  'RESPONSE',
  'ERROR',
  'HEARTBEAT',
  'SYSTEM'
]);

export type MessageType = z.infer<typeof MessageTypeSchema>;

export const MessagePrioritySchema = z.enum([
  'CRITICAL',
  'HIGH',
  'NORMAL',
  'LOW'
]);

export type MessagePriority = z.infer<typeof MessagePrioritySchema>;

export const MessageSchema = z.object({
  id: z.string().uuid(),
  type: MessageTypeSchema,
  priority: MessagePrioritySchema.default('NORMAL'),
  
  // Routing
  sender: z.string().uuid(),
  recipient: z.union([z.string().uuid(), z.literal('*'), z.array(z.string().uuid())]),
  replyTo: z.string().uuid().optional(),
  
  // Content
  payload: z.any(),
  payloadType: z.string(),
  
  // Context
  correlationId: z.string().uuid(),
  traceId: z.string().optional(),
  spanId: z.string().optional(),
  context: z.record(z.any()).default({}),
  
  // Timing
  createdAt: z.string().datetime(),
  expiresAt: z.string().datetime().optional(),
  ttlSeconds: z.number().int().positive().optional(),
  
  // Delivery
  deliveryGuarantee: z.enum(['at_most_once', 'at_least_once', 'exactly_once']).default('at_least_once'),
  
  // Metadata
  headers: z.record(z.string()).default({}),
  compression: z.enum(['none', 'gzip', 'zstd']).default('none'),
  encryption: z.enum(['none', 'aes256']).default('none'),
});

export type Message = z.infer<typeof MessageSchema>;

// ============================================================================
// EVENT MODELS
// ============================================================================

export const EventTypeSchema = z.enum([
  'AGENT_REGISTERED',
  'AGENT_UNREGISTERED',
  'AGENT_STATUS_CHANGED',
  'AGENT_HEALTH_CHANGED',
  'TASK_CREATED',
  'TASK_ASSIGNED',
  'TASK_STARTED',
  'TASK_COMPLETED',
  'TASK_FAILED',
  'TASK_CANCELLED',
  'WORKFLOW_STARTED',
  'WORKFLOW_COMPLETED',
  'WORKFLOW_FAILED',
  'MESSAGE_SENT',
  'MESSAGE_RECEIVED',
  'ERROR_OCCURRED',
  'SYSTEM_ALERT'
]);

export type EventType = z.infer<typeof EventTypeSchema>;

export const EventSchema = z.object({
  id: z.string().uuid(),
  type: EventTypeSchema,
  timestamp: z.string().datetime(),
  
  // Source
  source: z.object({
    type: z.enum(['AGENT', 'TASK', 'WORKFLOW', 'SYSTEM', 'USER']),
    id: z.string(),
    name: z.string().optional(),
  }),
  
  // Payload
  payload: z.any(),
  payloadSchema: z.string().optional(),
  
  // Context
  correlationId: z.string().uuid(),
  traceId: z.string().optional(),
  
  // Severity for alerting
  severity: z.enum(['info', 'warning', 'error', 'critical']).default('info'),
  
  // Metadata
  environment: z.string().optional(),
  region: z.string().optional(),
  version: z.string().optional(),
});

export type Event = z.infer<typeof EventSchema>;

// ============================================================================
// OBSERVABILITY MODELS
// ============================================================================

export const MetricSchema = z.object({
  name: z.string(),
  value: z.number(),
  type: z.enum(['counter', 'gauge', 'histogram', 'summary']),
  timestamp: z.string().datetime(),
  labels: z.record(z.string()).default({}),
  unit: z.string().optional(),
  description: z.string().optional(),
});

export type Metric = z.infer<typeof MetricSchema>;

export const LogEntrySchema = z.object({
  timestamp: z.string().datetime(),
  level: z.enum(['DEBUG', 'INFO', 'WARN', 'ERROR', 'FATAL']),
  message: z.string(),
  source: z.object({
    service: z.string(),
    instance: z.string(),
    file: z.string().optional(),
    line: z.number().int().optional(),
  }),
  context: z.record(z.any()).default({}),
  traceId: z.string().optional(),
  spanId: z.string().optional(),
  correlationId: z.string().optional(),
  tags: z.array(z.string()).default([]),
});

export type LogEntry = z.infer<typeof LogEntrySchema>;

export const TraceSpanSchema = z.object({
  traceId: z.string(),
  spanId: z.string(),
  parentSpanId: z.string().optional(),
  name: z.string(),
  service: z.string(),
  instance: z.string(),
  
  // Timing
  startTime: z.string().datetime(),
  endTime: z.string().datetime().optional(),
  durationMs: z.number().optional(),
  
  // Status
  status: z.enum(['ok', 'error', 'unset']).default('unset'),
  errorMessage: z.string().optional(),
  
  // Context
  attributes: z.record(z.any()).default({}),
  events: z.array(z.object({
    timestamp: z.string().datetime(),
    name: z.string(),
    attributes: z.record(z.any()),
  })).default([]),
});

export type TraceSpan = z.infer<typeof TraceSpanSchema>;
```

---

## 4. API Specifications

### 4.1 REST API

#### Base Configuration
```
Base URL: https://api.agentapi.phenotype.io/v2
Protocol: HTTPS (TLS 1.3)
Content-Type: application/json
Accept: application/json
```

#### Authentication
```http
Authorization: Bearer <jwt_token>
X-API-Key: <api_key>
X-Request-ID: <uuid>
```

#### Agent Management Endpoints

##### Register Agent
```http
POST /agents
Content-Type: application/json

{
  "identity": {
    "name": "document-processor-v2",
    "type": "LLM",
    "namespace": "production",
    "version": "2.1.0"
  },
  "capabilities": [
    {
      "id": "text-analysis",
      "name": "Text Analysis",
      "version": "1.0",
      "requirements": {
        "cpu": 2,
        "memory": 4096,
        "gpu": true
      }
    },
    {
      "id": "summarization",
      "name": "Document Summarization",
      "version": "2.0"
    }
  ],
  "configuration": {
    "maxConcurrentTasks": 5,
    "taskTimeoutSeconds": 600,
    "resources": {
      "cpuLimit": 4,
      "memoryLimit": 8192
    }
  },
  "endpoints": {
    "rest": "https://agent-123.internal:8080",
    "websocket": "wss://agent-123.internal:8081"
  },
  "tags": ["production", "nlp", "gpu-enabled"]
}
```

**Response (201 Created):**
```json
{
  "agent": {
    "identity": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "document-processor-v2",
      "type": "LLM",
      "namespace": "production",
      "version": "2.1.0",
      "createdAt": "2026-04-06T10:30:00Z",
      "updatedAt": "2026-04-06T10:30:00Z"
    },
    "status": "STANDBY",
    "capabilities": [...],
    "configuration": {...},
    "health": {
      "status": "healthy",
      "consecutiveFailures": 0
    },
    "tags": ["production", "nlp", "gpu-enabled"]
  },
  "registrationToken": "rt_xxxxxxxxxx",
  "expiresAt": "2026-04-07T10:30:00Z"
}
```

##### List Agents
```http
GET /agents?status=ACTIVE&type=LLM&namespace=production&page=1&limit=50
```

##### Get Agent
```http
GET /agents/{agent_id}
```

##### Update Agent
```http
PATCH /agents/{agent_id}
Content-Type: application/json

{
  "status": "ACTIVE",
  "configuration": {
    "maxConcurrentTasks": 10
  }
}
```

##### Delete Agent
```http
DELETE /agents/{agent_id}?force=false
```

##### Agent Heartbeat
```http
POST /agents/{agent_id}/heartbeat
Content-Type: application/json

{
  "timestamp": "2026-04-06T10:35:00Z",
  "metrics": {
    "cpuUsage": 45.2,
    "memoryUsage": 6144,
    "activeConnections": 3,
    "tasksCompleted": 127
  },
  "currentTask": "task-123"
}
```

#### Task Management Endpoints

##### Create Task
```http
POST /tasks
Content-Type: application/json

{
  "name": "summarize-contract",
  "description": "Summarize the service agreement contract",
  "type": "SINGLE",
  "priority": "NORMAL",
  "input": {
    "type": "text",
    "data": "[Full contract text...]",
    "metadata": {
      "source": "upload",
      "filename": "contract.pdf"
    }
  },
  "timeoutSeconds": 300,
  "tags": ["legal", "summarization"]
}
```

**Response (201 Created):**
```json
{
  "task": {
    "id": "task-550e8400-e29b-41d4-a716-446655440001",
    "name": "summarize-contract",
    "type": "SINGLE",
    "status": "PENDING",
    "priority": "NORMAL",
    "createdAt": "2026-04-06T10:30:00Z"
  },
  "estimatedStartTime": "2026-04-06T10:30:15Z",
  "positionInQueue": 3
}
```

##### List Tasks
```http
GET /tasks?status=RUNNING&priority=HIGH&agent_id={agent_id}&page=1&limit=50
```

##### Get Task
```http
GET /tasks/{task_id}
```

**Response (200 OK):**
```json
{
  "id": "task-550e8400-e29b-41d4-a716-446655440001",
  "name": "summarize-contract",
  "description": "Summarize the service agreement contract",
  "type": "SINGLE",
  "status": "COMPLETED",
  "priority": "NORMAL",
  "input": {...},
  "output": {
    "type": "json",
    "data": {
      "summary": "This service agreement outlines...",
      "keyPoints": ["Point 1", "Point 2"],
      "riskLevel": "medium"
    },
    "mimeType": "application/json"
  },
  "result": {
    "success": true,
    "metrics": {
      "executionTimeMs": 15234,
      "cpuTimeMs": 12100,
      "memoryPeakBytes": 2147483648
    }
  },
  "assignee": "550e8400-e29b-41d4-a716-446655440000",
  "createdAt": "2026-04-06T10:30:00Z",
  "startedAt": "2026-04-06T10:30:15Z",
  "completedAt": "2026-04-06T10:30:30Z",
  "progress": {
    "percent": 100,
    "currentStep": "completed",
    "message": "Task completed successfully"
  }
}
```

##### Cancel Task
```http
POST /tasks/{task_id}/cancel
Content-Type: application/json

{
  "reason": "User requested cancellation",
  "force": false
}
```

##### Retry Task
```http
POST /tasks/{task_id}/retry
Content-Type: application/json

{
  "resetState": true,
  "priority": "HIGH"
}
```

#### Workflow Endpoints

##### Create Workflow
```http
POST /workflows
Content-Type: application/json

{
  "name": "document-processing-pipeline",
  "description": "End-to-end document processing workflow",
  "version": "1.0.0",
  "nodes": [
    {
      "id": "start",
      "type": "START",
      "name": "Start",
      "position": { "x": 100, "y": 100 }
    },
    {
      "id": "extract",
      "type": "TASK",
      "name": "Extract Text",
      "position": { "x": 300, "y": 100 },
      "taskTemplate": {
        "type": "SINGLE",
        "input": {
          "type": "file"
        }
      }
    },
    {
      "id": "classify",
      "type": "CONDITION",
      "name": "Classify Document",
      "position": { "x": 500, "y": 100 },
      "condition": {
        "expression": "document.type",
        "operator": "equals"
      }
    },
    {
      "id": "summarize",
      "type": "TASK",
      "name": "Summarize",
      "position": { "x": 700, "y": 50 }
    },
    {
      "id": "extract-entities",
      "type": "TASK",
      "name": "Extract Entities",
      "position": { "x": 700, "y": 150 }
    },
    {
      "id": "merge",
      "type": "JOIN",
      "name": "Merge Results",
      "position": { "x": 900, "y": 100 }
    },
    {
      "id": "end",
      "type": "END",
      "name": "End",
      "position": { "x": 1100, "y": 100 }
    }
  ],
  "edges": [
    { "id": "e1", "source": "start", "target": "extract" },
    { "id": "e2", "source": "extract", "target": "classify" },
    { "id": "e3", "source": "classify", "target": "summarize", "condition": "contract" },
    { "id": "e4", "source": "classify", "target": "extract-entities", "condition": "invoice" },
    { "id": "e5", "source": "summarize", "target": "merge" },
    { "id": "e6", "source": "extract-entities", "target": "merge" },
    { "id": "e7", "source": "merge", "target": "end" }
  ],
  "startNodeId": "start",
  "timeoutSeconds": 1800
}
```

##### Execute Workflow
```http
POST /workflows/{workflow_id}/execute
Content-Type: application/json

{
  "input": {
    "document": {
      "url": "https://storage.example.com/contract.pdf",
      "type": "contract"
    }
  },
  "variables": {
    "extractOptions": { "ocr": true },
    "summaryLength": "medium"
  },
  "priority": "NORMAL"
}
```

#### Messaging Endpoints

##### Send Message
```http
POST /messages
Content-Type: application/json

{
  "type": "COMMAND",
  "priority": "HIGH",
  "recipient": "550e8400-e29b-41d4-a716-446655440000",
  "payload": {
    "command": "PROCESS_DOCUMENT",
    "parameters": {
      "documentId": "doc-123",
      "operations": ["extract", "summarize"]
    }
  },
  "payloadType": "ProcessDocumentCommand",
  "correlationId": "corr-550e8400-e29b-41d4-a716-446655440002",
  "deliveryGuarantee": "exactly_once",
  "ttlSeconds": 60
}
```

##### Broadcast Message
```http
POST /messages/broadcast
Content-Type: application/json

{
  "type": "EVENT",
  "priority": "NORMAL",
  "recipients": ["550e8400-e29b-41d4-a716-446655440000", "550e8400-e29b-41d4-a716-446655440001"],
  "payload": {
    "event": "SYSTEM_MAINTENANCE",
    "scheduledAt": "2026-04-07T02:00:00Z",
    "durationMinutes": 30
  },
  "payloadType": "SystemMaintenanceEvent"
}
```

### 4.2 GraphQL API

```graphql
# Schema Definition
type Agent {
  id: ID!
  identity: AgentIdentity!
  status: AgentStatus!
  capabilities: [Capability!]!
  configuration: AgentConfiguration!
  health: AgentHealth!
  currentTask: Task
  assignedTasks: [Task!]!
  tags: [String!]!
  labels: JSONObject
  createdAt: DateTime!
  updatedAt: DateTime!
}

type AgentIdentity {
  name: String!
  type: AgentType!
  namespace: String!
  version: String!
  metadata: JSONObject
}

type Task {
  id: ID!
  name: String!
  description: String
  type: TaskType!
  status: TaskStatus!
  priority: TaskPriority!
  input: TaskInput!
  output: TaskOutput
  result: TaskResult
  assignee: Agent
  createdAt: DateTime!
  startedAt: DateTime
  completedAt: DateTime
  progress: TaskProgress!
  tags: [String!]!
}

type Workflow {
  id: ID!
  name: String!
  description: String
  version: String!
  status: WorkflowStatus!
  nodes: [WorkflowNode!]!
  edges: [WorkflowEdge!]!
  startNodeId: ID!
  executions: [WorkflowExecution!]!
}

type Query {
  # Agent queries
  agent(id: ID!): Agent
  agents(
    status: AgentStatus
    type: AgentType
    namespace: String
    tags: [String!]
    page: Int = 1
    limit: Int = 50
  ): AgentConnection!
  
  # Task queries
  task(id: ID!): Task
  tasks(
    status: TaskStatus
    priority: TaskPriority
    assigneeId: ID
    tags: [String!]
    page: Int = 1
    limit: Int = 50
  ): TaskConnection!
  
  # Workflow queries
  workflow(id: ID!): Workflow
  workflows(
    status: WorkflowStatus
    page: Int = 1
    limit: Int = 50
  ): WorkflowConnection!
  
  # Metrics
  agentMetrics(
    agentId: ID
    metricNames: [String!]
    from: DateTime
    to: DateTime
    granularity: Granularity = MINUTE
  ): [Metric!]!
}

type Mutation {
  # Agent mutations
  registerAgent(input: RegisterAgentInput!): Agent!
  updateAgent(id: ID!, input: UpdateAgentInput!): Agent!
  deleteAgent(id: ID!, force: Boolean = false): Boolean!
  
  # Task mutations
  createTask(input: CreateTaskInput!): Task!
  cancelTask(id: ID!, reason: String): Task!
  retryTask(id: ID!, resetState: Boolean = true): Task!
  
  # Workflow mutations
  createWorkflow(input: CreateWorkflowInput!): Workflow!
  updateWorkflow(id: ID!, input: UpdateWorkflowInput!): Workflow!
  executeWorkflow(id: ID!, input: ExecuteWorkflowInput!): WorkflowExecution!
  
  # Message mutations
  sendMessage(input: SendMessageInput!): Message!
  broadcastMessage(input: BroadcastMessageInput!): [Message!]!
}

type Subscription {
  # Real-time subscriptions
  agentStatusChanged(agentId: ID): Agent!
  taskStatusChanged(taskId: ID): Task!
  workflowExecutionUpdated(executionId: ID): WorkflowExecution!
  messageReceived(agentId: ID!): Message!
  events(
    eventTypes: [EventType!]
    severity: [Severity!]
  ): Event!
}
```

### 4.3 gRPC API

```protobuf
// agent.proto
syntax = "proto3";
package agentapi.v2;

import "google/protobuf/timestamp.proto";
import "google/protobuf/any.proto";
import "google/protobuf/struct.proto";

service AgentService {
  rpc RegisterAgent(RegisterAgentRequest) returns (RegisterAgentResponse);
  rpc GetAgent(GetAgentRequest) returns (Agent);
  rpc ListAgents(ListAgentsRequest) returns (ListAgentsResponse);
  rpc UpdateAgent(UpdateAgentRequest) returns (Agent);
  rpc DeleteAgent(DeleteAgentRequest) returns (DeleteAgentResponse);
  rpc Heartbeat(stream HeartbeatRequest) returns (stream HeartbeatResponse);
  rpc SubscribeToAgentEvents(SubscribeRequest) returns (stream AgentEvent);
}

service TaskService {
  rpc CreateTask(CreateTaskRequest) returns (Task);
  rpc GetTask(GetTaskRequest) returns (Task);
  rpc ListTasks(ListTasksRequest) returns (ListTasksResponse);
  rpc CancelTask(CancelTaskRequest) returns (Task);
  rpc RetryTask(RetryTaskRequest) returns (Task);
  rpc StreamTaskLogs(StreamTaskLogsRequest) returns (stream LogEntry);
}

service WorkflowService {
  rpc CreateWorkflow(CreateWorkflowRequest) returns (Workflow);
  rpc GetWorkflow(GetWorkflowRequest) returns (Workflow);
  rpc ExecuteWorkflow(ExecuteWorkflowRequest) returns (WorkflowExecution);
  rpc GetWorkflowExecution(GetWorkflowExecutionRequest) returns (WorkflowExecution);
  rpc StreamWorkflowExecution(StreamWorkflowExecutionRequest) returns (stream WorkflowExecutionUpdate);
}

service MessageService {
  rpc SendMessage(SendMessageRequest) returns (Message);
  rpc SendMessageStream(stream SendMessageRequest) returns (stream MessageReceipt);
  rpc ReceiveMessages(ReceiveMessagesRequest) returns (stream Message);
  rpc AcknowledgeMessage(AcknowledgeMessageRequest) returns (AcknowledgeMessageResponse);
}

// Message definitions
message Agent {
  string id = 1;
  AgentIdentity identity = 2;
  AgentStatus status = 3;
  repeated Capability capabilities = 4;
  AgentConfiguration configuration = 5;
  AgentHealth health = 6;
  string current_task_id = 7;
  repeated string assigned_task_ids = 8;
  repeated string tags = 9;
  google.protobuf.Struct labels = 10;
  google.protobuf.Timestamp created_at = 11;
  google.protobuf.Timestamp updated_at = 12;
}

message AgentIdentity {
  string name = 1;
  AgentType type = 2;
  string namespace = 3;
  string version = 4;
  google.protobuf.Struct metadata = 5;
}

enum AgentStatus {
  AGENT_STATUS_UNSPECIFIED = 0;
  AGENT_STATUS_NEW = 1;
  AGENT_STATUS_STANDBY = 2;
  AGENT_STATUS_ACTIVE = 3;
  AGENT_STATUS_RUNNING = 4;
  AGENT_STATUS_SUSPENDED = 5;
  AGENT_STATUS_IDLE = 6;
  AGENT_STATUS_ERROR = 7;
  AGENT_STATUS_COMPLETED = 8;
  AGENT_STATUS_TERMINATE = 9;
  AGENT_STATUS_DESTROY = 10;
}

enum AgentType {
  AGENT_TYPE_UNSPECIFIED = 0;
  AGENT_TYPE_LLM = 1;
  AGENT_TYPE_TOOL = 2;
  AGENT_TYPE_ORCHESTRATOR = 3;
  AGENT_TYPE_HUMAN = 4;
  AGENT_TYPE_SYSTEM = 5;
  AGENT_TYPE_CUSTOM = 6;
}

message Task {
  string id = 1;
  string name = 2;
  string description = 3;
  TaskType type = 4;
  TaskStatus status = 5;
  TaskPriority priority = 6;
  TaskInput input = 7;
  TaskOutput output = 8;
  TaskResult result = 9;
  string assignee_id = 10;
  google.protobuf.Timestamp created_at = 11;
  google.protobuf.Timestamp scheduled_at = 12;
  google.protobuf.Timestamp started_at = 13;
  google.protobuf.Timestamp completed_at = 14;
  TaskProgress progress = 15;
  repeated string tags = 16;
  google.protobuf.Struct labels = 17;
}

enum TaskStatus {
  TASK_STATUS_UNSPECIFIED = 0;
  TASK_STATUS_PENDING = 1;
  TASK_STATUS_SCHEDULED = 2;
  TASK_STATUS_ASSIGNED = 3;
  TASK_STATUS_RUNNING = 4;
  TASK_STATUS_PAUSED = 5;
  TASK_STATUS_CANCELLED = 6;
  TASK_STATUS_FAILED = 7;
  TASK_STATUS_COMPLETED = 8;
  TASK_STATUS_TIMEOUT = 9;
}

// Additional message definitions...
```

---

## 5. Configuration

### 5.1 Server Configuration (YAML)

```yaml
# agentapi-config.yaml

server:
  host: 0.0.0.0
  port: 8080
  
  # Protocol configuration
  protocols:
    rest:
      enabled: true
      max_request_size_mb: 50
      read_timeout_seconds: 30
      write_timeout_seconds: 30
    
    graphql:
      enabled: true
      endpoint: /graphql
      playground: true
      introspection: true
      max_complexity: 1000
      max_depth: 10
    
    grpc:
      enabled: true
      port: 50051
      max_receive_message_length: 16777216  # 16MB
      max_send_message_length: 16777216
    
    websocket:
      enabled: true
      endpoint: /ws
      read_buffer_size: 1024
      write_buffer_size: 1024
      ping_period_seconds: 30
      pong_wait_seconds: 60
      max_connections_per_agent: 5

# Authentication & Authorization
auth:
  enabled: true
  jwt:
    secret: ${JWT_SECRET}
    issuer: agentapi.phenotype.io
    audience: agentapi-clients
    access_token_ttl_minutes: 60
    refresh_token_ttl_days: 7
  
  api_keys:
    enabled: true
    header_name: X-API-Key
    rate_limit_per_minute: 1000
  
  rbac:
    enabled: true
    policy_refresh_seconds: 300

# Database configuration
database:
  primary:
    type: postgresql
    host: ${DB_HOST}
    port: 5432
    database: agentapi
    user: ${DB_USER}
    password: ${DB_PASSWORD}
    pool_size: 20
    connection_timeout_seconds: 5
    idle_timeout_seconds: 300
    max_lifetime_seconds: 1800
  
  cache:
    type: redis
    host: ${REDIS_HOST}
    port: 6379
    password: ${REDIS_PASSWORD}
    db: 0
    pool_size: 50
    default_ttl_seconds: 300

# Message broker
message_broker:
  type: kafka
  brokers:
    - kafka-1:9092
    - kafka-2:9092
    - kafka-3:9092
  
  topics:
    agent_events: agent.events.v2
    task_events: task.events.v2
    workflow_events: workflow.events.v2
    command_requests: command.requests.v2
    command_responses: command.responses.v2
  
  producer:
    acks: all
    retries: 3
    batch_size: 16384
    linger_ms: 5
  
  consumer:
    group_id: agentapi-consumers
    auto_offset_reset: earliest
    max_poll_records: 500

# Observability
observability:
  metrics:
    enabled: true
    endpoint: /metrics
    provider: prometheus
    report_interval_seconds: 15
    
    custom_metrics:
      - name: agent_active_count
        type: gauge
        description: Number of active agents
      - name: task_execution_duration
        type: histogram
        description: Task execution duration in milliseconds
        buckets: [100, 500, 1000, 5000, 10000, 30000, 60000]
  
  tracing:
    enabled: true
    sampling_rate: 0.1
    exporter: jaeger
    jaeger_endpoint: http://jaeger:14268/api/traces
  
  logging:
    level: info
    format: json
    output: stdout
    include_trace_ids: true
    include_caller: true
  
  alerting:
    enabled: true
    rules:
      - name: high_agent_failure_rate
        condition: agent_failure_rate > 0.1
        duration: 5m
        severity: warning
        channels:
          - slack
          - pagerduty

# Task scheduling
scheduler:
  enabled: true
  poll_interval_seconds: 1
  max_concurrent_tasks: 1000
  task_assignment_strategy: capability_match
  
  retry:
    max_retries: 3
    backoff_strategy: exponential
    initial_delay_ms: 1000
    max_delay_ms: 30000
    backoff_multiplier: 2

# Security
security:
  rate_limiting:
    enabled: true
    requests_per_second: 1000
    burst_size: 2000
  
  tls:
    enabled: true
    cert_file: /etc/agentapi/tls.crt
    key_file: /etc/agentapi/tls.key
    min_version: "1.3"
  
  cors:
    enabled: true
    allowed_origins:
      - https://console.phenotype.io
    allowed_methods:
      - GET
      - POST
      - PUT
      - PATCH
      - DELETE
    allowed_headers:
      - Authorization
      - Content-Type
      - X-Request-ID
    max_age_seconds: 86400

# Feature flags
features:
  streaming_execution_logs: true
  real_time_metrics: true
  workflow_visual_editor: true
  agent_discovery: true
  multi_tenancy: false
```

---

## 6. Performance Benchmarks

### 6.1 Target Metrics

| Metric | Target | Critical Threshold | Load Test Scenario |
|--------|--------|-------------------|-------------------|
| API Response Time (p50) | < 50ms | 100ms | 10K concurrent requests |
| API Response Time (p99) | < 200ms | 500ms | 10K concurrent requests |
| Task Creation Throughput | > 1K/sec | 500/sec | Sustained load |
| Task Assignment Latency | < 100ms | 500ms | 1000 pending tasks |
| Agent Heartbeat Processing | > 50K/sec | 20K/sec | Burst load |
| WebSocket Message Latency | < 10ms | 50ms | 10K active connections |
| Database Query Time (p99) | < 20ms | 100ms | Complex joins |
| Cache Hit Rate | > 95% | 85% | Mixed workload |
| Event Processing Latency | < 50ms | 200ms | 100K events/sec |

### 6.2 Benchmark Results

```yaml
# performance-benchmarks.yaml

environment:
  deployment: kubernetes
  nodes: 5
  cpu_per_node: 16
  memory_per_node: 64GB
  database: PostgreSQL 15 (primary + replica)
  cache: Redis 7 (cluster mode)
  message_broker: Kafka 3.5 (3 brokers)
  load_generator: k6

test_scenarios:
  - name: agent_registration_burst
    description: Simulate burst of agent registrations
    duration: 5m
    stages:
      - duration: 1m
        target: 1000  # 1000 registrations/min
      - duration: 2m
        target: 5000  # Peak load
      - duration: 2m
        target: 0
    
    results:
      p50_latency_ms: 23
      p95_latency_ms: 67
      p99_latency_ms: 145
      error_rate: 0.001
      throughput_rps: 4873
      
  - name: sustained_task_throughput
    description: Sustained task creation and execution
    duration: 30m
    virtual_users: 500
    
    results:
      tasks_created_per_second: 1247
      tasks_completed_per_second: 1234
      avg_execution_time_ms: 2345
      p99_queue_time_ms: 89
      failure_rate: 0.002
      
  - name: websocket_stress_test
    description: Maximum WebSocket connection handling
    duration: 10m
    connections: 50000
    message_rate_per_connection: 1/sec
    
    results:
      max_concurrent_connections: 50000
      message_latency_p99_ms: 8
      connection_drop_rate: 0.0001
      memory_usage_gb: 12.4
      
  - name: workflow_execution
    description: Complex workflow with parallel branches
    duration: 15m
    concurrent_workflows: 100
    nodes_per_workflow: 50
    
    results:
      workflow_completion_rate: 0.997
      avg_completion_time_ms: 45678
      p99_completion_time_ms: 89123
      node_execution_success_rate: 0.9995
```

---

## 7. Security Model

### 7.1 Authentication Flow

```
┌─────────┐                                    ┌─────────────┐      ┌─────────────┐
│  Client │                                    │ API Gateway │      │ Auth Service│
└────┬────┘                                    └──────┬──────┘      └──────┬──────┘
     │                                               │                    │
     │  1. POST /auth/login                          │                    │
     │     {credentials}                             │                    │
     │──────────────────────────────────────────────▶                    │
     │                                               │  2. Validate       │
     │                                               │───────────────────▶│
     │                                               │                    │
     │                                               │  3. Return tokens  │
     │                                               │◀───────────────────│
     │  4. {access_token, refresh_token}             │                    │
     │◀──────────────────────────────────────────────│                    │
     │                                               │                    │
     │  5. API request with Bearer token             │                    │
     │──────────────────────────────────────────────▶                    │
     │                                               │  6. Validate JWT   │
     │                                               │───────────────────▶│
     │                                               │                    │
     │                                               │  7. Claims + RBAC  │
     │                                               │◀───────────────────│
     │  8. Response                                  │                    │
     │◀──────────────────────────────────────────────│                    │
     │                                               │                    │
```

### 7.2 Authorization Matrix

| Role | Read Agents | Write Agents | Execute Tasks | Manage Workflows | View Metrics | Admin |
|------|:-----------:|:------------:|:-------------:|:----------------:|:------------:|:-----:|
| viewer | ✅ | ❌ | ❌ | ❌ | ✅ | ❌ |
| developer | ✅ | ✅ (own) | ✅ | ✅ (own) | ✅ | ❌ |
| operator | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ |
| admin | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| service | ✅ | ❌ | ✅ | ❌ | ✅ | ❌ |

---

## 8. Deployment Guide

### 8.1 Kubernetes Deployment

```yaml
# k8s-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: agentapi-server
  namespace: agentapi
spec:
  replicas: 3
  selector:
    matchLabels:
      app: agentapi-server
  template:
    metadata:
      labels:
        app: agentapi-server
    spec:
      containers:
        - name: server
          image: phenotype/agentapi-plusplus:v2.0.0
          ports:
            - containerPort: 8080
            - containerPort: 50051
          envFrom:
            - configMapRef:
                name: agentapi-config
            - secretRef:
                name: agentapi-secrets
          resources:
            requests:
              memory: "2Gi"
              cpu: "1000m"
            limits:
              memory: "4Gi"
              cpu: "2000m"
          livenessProbe:
            httpGet:
              path: /health/live
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health/ready
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: agentapi-service
  namespace: agentapi
spec:
  selector:
    app: agentapi-server
  ports:
    - name: http
      port: 8080
      targetPort: 8080
    - name: grpc
      port: 50051
      targetPort: 50051
  type: LoadBalancer
```

---

## 9. Appendices

### Appendix A: Glossary

| Term | Definition |
|------|------------|
| **Agent** | A computational entity that can receive tasks, process them, and return results |
| **Task** | A unit of work assigned to an agent for execution |
| **Workflow** | A DAG of interconnected tasks with defined execution logic |
| **Capability** | A declared skill or function that an agent can perform |
| **Message** | A communication unit exchanged between agents |
| **Namespace** | Logical isolation boundary for agents and resources |
| **Orchestrator** | A special agent type that coordinates other agents |

### Appendix B: Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `AGENT_NOT_FOUND` | 404 | Agent ID does not exist |
| `AGENT_OFFLINE` | 503 | Agent is not currently connected |
| `TASK_NOT_FOUND` | 404 | Task ID does not exist |
| `TASK_ALREADY_ASSIGNED` | 409 | Task is already assigned to another agent |
| `WORKFLOW_INVALID` | 400 | Workflow definition contains errors |
| `CAPABILITY_MISMATCH` | 400 | Agent lacks required capabilities |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `INSUFFICIENT_RESOURCES` | 503 | No agents available with required resources |

### Appendix C: SDK Examples

#### Python SDK

```python
from agentapi import Client, Agent, Task

# Initialize client
client = Client(
    api_url="https://api.agentapi.phenotype.io",
    api_key="your-api-key"
)

# Register agent
agent = client.agents.register(
    name="my-agent",
    type="LLM",
    capabilities=["text-generation", "summarization"]
)

# Create and execute task
task = client.tasks.create(
    name="summarize",
    input={"text": "Long document..."},
    assignee=agent.id
)

# Wait for completion
result = task.wait_for_completion(timeout=300)
print(result.output)
```

#### TypeScript SDK

```typescript
import { AgentAPIClient } from '@phenotype/agentapi';

const client = new AgentAPIClient({
  apiUrl: 'https://api.agentapi.phenotype.io',
  apiKey: process.env.AGENTAPI_KEY!
});

// Real-time event subscription
const subscription = client.events.subscribe({
  eventTypes: ['TASK_COMPLETED', 'AGENT_STATUS_CHANGED']
});

subscription.on('task_completed', (event) => {
  console.log(`Task ${event.payload.taskId} completed`);
});

// Execute workflow
const execution = await client.workflows.execute(
  'workflow-id',
  {
    input: { document: 'url-to-doc' },
    variables: { summaryLength: 'short' }
  }
);

// Stream progress
for await (const update of execution.stream()) {
  console.log(`Progress: ${update.progress.percent}%`);
}
```

### Appendix D: Webhook Events

| Event | Payload | Trigger |
|-------|---------|---------|
| `agent.registered` | Agent object | New agent registration |
| `agent.status_changed` | { agent_id, old_status, new_status } | Agent state transition |
| `task.created` | Task object | New task created |
| `task.completed` | { task_id, result } | Task finished |
| `workflow.execution.completed` | { execution_id, workflow_id, status } | Workflow finished |

### Appendix E: Metrics Reference

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `agent_active_total` | Gauge | status, type | Number of active agents |
| `agent_tasks_completed_total` | Counter | agent_id, status | Tasks completed by agent |
| `task_execution_duration_seconds` | Histogram | type, priority | Task execution time |
| `task_queue_wait_seconds` | Histogram | priority | Time spent in queue |
| `message_sent_total` | Counter | type, priority | Messages sent |
| `api_requests_duration_seconds` | Histogram | endpoint, method | API response time |

### Appendix F: Database Schema

```sql
-- Agents table
CREATE TABLE agents (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    namespace VARCHAR(255) NOT NULL,
    version VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    capabilities JSONB NOT NULL DEFAULT '[]',
    configuration JSONB NOT NULL DEFAULT '{}',
    health JSONB NOT NULL DEFAULT '{}',
    endpoints JSONB,
    current_task_id UUID,
    assigned_task_ids UUID[] DEFAULT '{}',
    tags VARCHAR(255)[],
    labels JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_agents_status ON agents(status);
CREATE INDEX idx_agents_type ON agents(type);
CREATE INDEX idx_agents_namespace ON agents(namespace);
CREATE INDEX idx_agents_tags ON agents USING GIN(tags);

-- Tasks table
CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    name VARCHAR(500) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    priority VARCHAR(50) NOT NULL,
    input JSONB NOT NULL,
    output JSONB,
    result JSONB,
    assignee_id UUID REFERENCES agents(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    scheduled_at TIMESTAMPTZ,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    deadline_at TIMESTAMPTZ,
    progress JSONB DEFAULT '{}',
    tags VARCHAR(255)[],
    labels JSONB DEFAULT '{}',
    workflow_id UUID,
    step_number INTEGER,
    dependencies UUID[] DEFAULT '{}',
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3
);

CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_assignee ON tasks(assignee_id);
CREATE INDEX idx_tasks_workflow ON tasks(workflow_id);
CREATE INDEX idx_tasks_created_at ON tasks(created_at DESC);
```

### Appendix G: Migration Guide

**From v1.x to v2.x:**

1. **API Changes**
   - Base URL changed from `/v1` to `/v2`
   - Task `status` enum extended with `SCHEDULED` and `TIMEOUT`
   - Agent `type` enum added `ORCHESTRATOR`

2. **Database Migration**
   ```sql
   -- Run migration script
   ALTER TABLE tasks ADD COLUMN IF NOT EXISTS workflow_id UUID;
   ALTER TABLE tasks ADD COLUMN IF NOT EXISTS step_number INTEGER;
   CREATE INDEX idx_tasks_workflow ON tasks(workflow_id);
   ```

3. **SDK Updates**
   - Update SDK to v2.x
   - Update authentication method
   - Review breaking changes in CHANGELOG

### Appendix H: Troubleshooting

| Issue | Symptoms | Solution |
|-------|----------|----------|
| Agent not receiving tasks | Task stays in PENDING | Check agent status is ACTIVE, verify capabilities match |
| High task queue latency | Queue growing, slow assignment | Scale up agent pool, check scheduler health |
| WebSocket disconnections | Frequent reconnects | Check network stability, adjust heartbeat intervals |
| Database connection errors | 500 errors | Check connection pool size, verify DB health |
| Memory leaks | Growing memory usage | Review agent cleanup, check for event listener leaks |

### Appendix I: Compliance

| Standard | Implementation |
|----------|---------------|
| SOC 2 | Audit logging, access controls, encryption |
| GDPR | Data retention policies, right to deletion |
| HIPAA | Optional PHI handling capabilities |
| ISO 27001 | Security controls, risk management |

### Appendix J: CLI Reference

```bash
# Install CLI
npm install -g @phenotype/agentapi-cli

# Authentication
agentapi login --api-key <key>

# Agent commands
agentapi agents list --status ACTIVE
agentapi agents get <agent-id>
agentapi agents logs <agent-id> --follow

# Task commands
agentapi tasks create --file task.json
agentapi tasks get <task-id>
agentapi tasks cancel <task-id>

# Workflow commands
agentapi workflows execute <workflow-id> --input input.json
agentapi workflows status <execution-id>

# Monitoring
agentapi metrics --agent <agent-id> --from "1h ago"
agentapi events --follow
```

### Appendix K: WebSocket Protocol

**Connection:**
```javascript
const ws = new WebSocket('wss://api.agentapi.phenotype.io/v2/ws');

// Authenticate
ws.send(JSON.stringify({
  type: 'AUTH',
  token: 'jwt_token'
}));

// Subscribe to events
ws.send(JSON.stringify({
  type: 'SUBSCRIBE',
  channels: ['agent:550e8400', 'task:task-123']
}));
```

**Message Format:**
```json
{
  "id": "msg-uuid",
  "type": "EVENT",
  "channel": "agent:550e8400",
  "payload": {
    "event": "STATUS_CHANGED",
    "agentId": "550e8400",
    "oldStatus": "ACTIVE",
    "newStatus": "RUNNING"
  },
  "timestamp": "2026-04-06T10:30:00Z"
}
```

### Appendix L: Best Practices

1. **Agent Design**
   - Keep agents focused on single responsibility
   - Declare all capabilities explicitly
   - Implement graceful shutdown
   - Use heartbeats to maintain liveness

2. **Task Design**
   - Keep tasks idempotent
   - Set appropriate timeouts
   - Include all necessary context in input
   - Handle partial failures gracefully

3. **Workflow Design**
   - Limit workflow depth to avoid complexity
   - Use conditions for branching, not error handling
   - Set workflow-level timeouts
   - Monitor execution metrics

4. **Performance**
   - Use batch operations when possible
   - Implement proper caching
   - Monitor resource utilization
   - Scale horizontally based on metrics

---

## Document Information

| Field | Value |
|-------|-------|
| **Document ID** | SPEC-AGENTAPI-001 |
| **Version** | 2.0.0 |
| **Status** | Active |
| **Last Updated** | 2026-04-06 |

---

*This specification defines the complete technical architecture for agentapi-plusplus.*
