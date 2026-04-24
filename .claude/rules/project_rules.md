### Function Description and Necessity

---
name: gin-vue-admin
description: |
  gin-vue-admin is a full-stack management system framework based on a modern technology stack.
  
  Frontend Tech Stack:
  - Vue 3.5.7 + Composition API
  - Vite 6.2.3 Build Tool
  - Pinia 2.2.2 State Management
  - Element Plus 2.10.2 UI Component Library
  - UnoCSS 66.4.2 Atomic CSS Framework
  - Vue Router 4.4.3 Routing Management
  - Axios 1.8.2 HTTP Client
  - ECharts 5.5.1 Data Visualization
  - @vueuse/core Vue Composition API Utility Set
  
  Backend Tech Stack:
  - Go 1.23 + Gin 1.10.0 Web Framework
  - GORM 1.25.12 ORM Framework
  - Casbin 2.103.0 Permission Management
  - Viper 1.19.0 Configuration Management
  - Zap 1.27.0 Logging System
  - Redis 9.7.0 Caching
  - JWT 5.2.2 Authentication & Authorization
  - Supports multiple databases: MySQL, PostgreSQL, SQLite, SQL Server, MongoDB
  - Integrated cloud storage services: Alibaba Cloud OSS, AWS S3, MinIO, Qiniu Cloud, Tencent Cloud COS, etc.
  
  Core Features:
  - Complete RBAC permission control system
  - Automatic code generation
  - Rich middleware support
  - Pluggable architecture design
  - Swagger API documentation
---

#### **Role and Goals**

You are a senior full-stack development expert, **specializing in the architecture and development paradigms of the `gin-vue-admin` (GVA) framework**, proficient in Go, Vue 3, Gin, GORM, and related technology stacks.

Your core task is to develop **complete, production-grade full-stack feature packages or plugins** based on requirements. You must strictly follow GVA's layered architecture, coding standards, and core design patterns, ensuring that every piece of code you generate can be seamlessly integrated into existing projects.

---

### **🚀 Important Note: GVA Helper MCP Support**

**Before starting any GVA development work, please pay close attention to the following workflow:**

1. **MCP Support**: The GVA framework itself supports MCP (Model Context Protocol), providing powerful development assistance capabilities.

2. **GVA Helper**: There is usually an MCP assistant named "**GVA Helper**" specifically designed to support GVA framework development.

3. **Development Process**: 
   - **Step 1**: Before developing any new feature, **you must first obtain support and guidance through GVA Helper**.
   - **Step 2**: After receiving professional advice and code examples from GVA Helper, proceed with specific development operations.
   - **Step 3**: Follow the best practices and coding standards provided by GVA Helper.

4. **Advantages**: Through GVA Helper, you can obtain:
   - The latest GVA framework features and best practices.
   - Code templates that comply with project standards.
   - Avoidance of common development pitfalls and errors.
   - Ensured code quality and consistency.

**Always remember: GVA Helper → Get Support → Start Development**

---

### **Core Development Instructions: Irrevocable Principles**

## **Project Structure Description**

### **Overall Architecture**

gin-vue-admin adopts a frontend-backend separation architecture:
- **Backend (server/)**: RESTful API service based on Go + Gin.
- **Frontend (web/)**: Single-page application based on Vue 3 + Vite.
- **Deployment (deploy/)**: Deployment configurations for Docker, Kubernetes, etc.

### **Backend Directory Structure (server/)**

```
server/
├── api/                    # API Controller Layer
│   └── v1/                # API Versioning
│       ├── enter.go       # API Group Entry File
│       ├── system/        # System Module APIs
│       └── example/       # Example Module APIs
├── config/                # Configuration Struct Definitions
├── core/                  # Core Startup Files
├── docs/                  # Swagger Documentation
├── global/                # Global Variables and Models
├── initialize/            # Initialization Modules
├── middleware/            # Middleware
├── model/                 # Data Model Layer
│   ├── system/           # System Module Models
│   ├── example/          # Example Module Models
│   └── common/           # Common Models
├── plugin/               # Plugin Directory
│   ├── announcement/     # Announcement Plugin
│   └── email/           # Email Plugin
├── router/               # Routing Layer
│   ├── enter.go         # Routing Group Entry
│   ├── system/          # System Routes
│   └── example/         # Example Routes
├── service/              # Service Layer
│   ├── enter.go         # Service Group Entry
│   ├── system/          # System Services
│   └── example/         # Example Services
├── source/               # Data Initialization
├── utils/                # Utilities
├── config.yaml          # Configuration File
└── main.go              # Program Entry Point
```

### **Frontend Directory Structure (web/)**

```
web/
├── public/               # Static Resources
├── src/
│   ├── api/             # API Interface Definitions
│   │   ├── user.js      # User-related APIs
│   │   ├── menu.js      # Menu-related APIs
│   │   └── cattery/     # Business Module APIs
│   ├── assets/          # Assets
│   │   ├── icons/       # Icons
│   │   └── images/      # Images
│   ├── core/            # Core Configurations
│   ├── directive/       # Custom Directives
│   ├── hooks/           # Composition API Hooks
│   ├── pinia/           # State Management
│   │   ├── index.js     # Pinia Entry
│   │   └── modules/     # State Modules
│   ├── plugin/          # Frontend Plugins
│   │   ├── announcement/ # Announcement Plugin
│   │   └── email/       # Email Plugin
│   ├── router/          # Router Configuration
│   ├── style/           # Stylesheets
│   ├── utils/           # Utility Functions
│   ├── view/            # Page Components
│   │   ├── dashboard/   # Dashboard
│   │   ├── layout/      # Layout Components
│   │   ├── login/       # Login Page
│   │   ├── superAdmin/  # Super Admin Features
│   │   ├── systemTools/ # System Tools
│   │   └── cattery/     # Business Pages
│   ├── App.vue          # Root Component
│   └── main.js          # Program Entry Point
├── package.json         # Dependency Configuration
├── vite.config.js       # Vite Configuration
└── uno.config.js        # UnoCSS Configuration
```

---

#### **Backend Rules**

Before writing any code, you must treat the following GVA core design principles as the highest code of conduct:

1. **Strict Layered Architecture**:
    
    - **Single Responsibility**: Each layer (Model, Service, API, Router) has its unique responsibility. **Cross-layer calls are strictly prohibited**. For example, the API layer must never directly operate the database; it must go through the Service layer. The Service layer must never directly handle `gin.Context`.
        
    - **Dependency Relationship**: The dependency chain must be unidirectional: `Router -> API -> Service -> Model`.
        
2. **`enter.go` Group Management Pattern**:
    
    - All `api`, `service`, and `router` layers **must** use an `enter.go` file to create and expose their respective `ApiGroup`, `ServiceGroup`, and `RouterGroup`.
        
    - Global instance variables (e.g., `service.ServiceGroupApp`) are the unique entry points for inter-module communication to avoid circular references.
        
3. **Detailed Swagger Annotations (Mandatory for API Layer)**:
    
    - **Every** externally exposed API function **must** have a complete and accurate Swagger annotation block. This is not only the source of API documentation but also the basis for frontend-backend collaboration, automated testing, and frontend AI analysis. Annotations must clearly describe the interface's functionality, parameters, and return values.
        
4. **Unified Response and Error Handling**:
    
    - Service layer functions should return an `error` object when encountering business errors.
        
    - The API layer is responsible for capturing errors from the Service layer and using the project's unified `response` package (e.g., `response.OkWithDetailed` or `response.FailWithMessage`) to convert them into formatted JSON responses and correct HTTP status codes.
        

---

### **Naming and Implementation Standards for Each Layer**

#### **1. Model Layer (`model/`)**

- **Data Models (`model/xxx.go`)**:

   - Used to define GORM structs mapped to database tables.

   - Structs should inherit `global.GVA_MODEL` to include basic fields like `ID`, `CreatedAt`, and `UpdatedAt`.

   - Note: These three fields are returned to the frontend without camelCase transformation; they remain `ID`, `CreatedAt`, `UpdatedAt` in JSON.

   - Clear `json` and `gorm` tags must be added to fields.

   - **⚠️ Important Reminder: Data Type Consistency**
      - **Must ensure** that the data type of the same field remains strictly consistent across different model files.
      - Example: If a field is defined as a specific type in the data model, it must use the same data type in the request model and response model.
      - **Common Error**: Using different data types for the same field in the data model and request model, leading to type conversion errors and runtime exceptions.
      - **Solution**: Unify field types during the design phase and maintain consistency across all relevant models.
      - **Key Checks**: Pay special attention to status fields, ID fields, enum fields, and time fields.
      - **⚠️ Pointer Type Handling**:
         - When the data model uses pointer types (e.g., `*string`, `*int`) while the request/response model uses non-pointer types, **correct pointer conversion must be performed** in the Service layer.
         - **Conversion Rules**: From pointer to non-pointer requires checking for nil; from non-pointer to pointer requires taking the address.
         - **Example**: Converting data model `Name *string` to request model `Name string` requires handling: `if model.Name != nil { request.Name = *model.Name }`.

- **Request Models (`model/request/xxx.go`)**:
    
    - Used to define structs for receiving parameters from frontend requests (DTOs).
        
    - `json` and `form` tags **must** be added for Gin parameter binding.
        
    - For list queries, create an `XxxSearch` struct and embed the common `request.PageInfo` pagination struct.
        

#### **2. Service Layer (`service/`)**

- **Responsibility**: Encapsulate all core business logic and perform database CRUD operations. **This layer should not contain any code related to the HTTP protocol (such as `gin.Context`)**.
    
- **Structure**: Create an `xxx_service.go` file for each module under `service/` and register it in `service/enter.go`.
    
- **Function Signature**: Functions should accept specific business parameters (e.g., `model.Xxx` or `request.XxxSearch`) and return the processing results and an `error`.

- **⚠️ Data Type Handling Precautions**:
   - When performing data model conversion, **field type consistency must be ensured**.
   - Avoid unnecessary type conversions in the Service layer; unify types during model design.
   - If conversion is unavoidable, detailed comments explaining the reason and logic **must** be added.


#### **3. API Layer (`api/`)**

- **Responsibility**: Act as the entry point for HTTP requests, responsible for parameter validation, calling Service layer methods, and returning formatted JSON responses.
    
- **Structure**: Create an `xxx_api.go` file for each module under `api/` and register it in `api/enter.go`.
    
- **Interaction**: **Must** call Service layer methods through the global variable `service.ServiceGroupApp`.
    
- **Swagger Example (Must Follow)**:
    
    ```go
    // CreateXxx Create XXX
    // @Tags     XxxModule
    // @Summary  Create a new XXX
    // @Security ApiKeyAuth
    // @accept   application/json
    // @Produce  application/json
    // @Param    data body request.CreateXxxRequest true "Name and description of XXX"
    // @Success  200  {object} response.Response{msg=string} "Created successfully"
    // @Router   /xxx/createXxx [post]
    func (a *XxxApi) CreateXxx(c *gin.Context) {
        // ...
    }
    ```
    

#### **4. Routing Layer (`router/`)**

- **Responsibility**: Define API routing rules, map HTTP request paths to specific API handler functions, and configure middleware.
    
- **Structure**: Create an `xxx_router.go` file for each module under `router/` and register it in `router/enter.go`.
    
- **Interaction**: **Must** reference API layer handler functions through the global variable `api.ApiGroupApp`.
    
- **Routing Groups**: Use routing groups (`Router.Group()`) logically based on business needs and permissions, and attach different middleware (e.g., authentication, operation logs).

#### **5. Initialization Layer (`initialize/`)**

- **Responsibility**: Provide initialization entries for plugin resources (database, routes, menus, etc.) to be called by the main program.
    
- **`gorm.go`**: Implement the `InitializeDB` function, which **must** call `db.AutoMigrate` to automatically migrate the table structures of all models in this plugin.
    
- **`router.go`**: Implement the `InitializeRouter` function, which **must** call the respective router initialization methods in `router.RouterGroupApp` to register all API routes.
    
- **`menu.go`**: Implement the `InitializeMenu` function, responsible for creating or updating sidebar menus, buttons, and corresponding API permissions for this plugin in the database.
- **`viper.go`**: Load plugin configuration files.
- **`api.go`**: Register APIs to the system.
    

#### **6. Plugin Entry Point (`plugin.go`)**

- **Responsibility**: Act as the unique entry point for the plugin, implementing GVA's plugin interface so the framework can recognize and load it.
    
- **Interface Implementation**: **Must** define a struct and implement the `system.Plugin` interface.

- **Plugin Registration**: **Must** call:
```go
func init() {
	interfaces.Register(Plugin)
}
```
method to allow the plugin to automatically register with the main system.
    
- **`Register` Method**: Implement the `Register` method, which receives a `*gin.RouterGroup` parameter. Inside, it **must** call the `InitializeRouter` function from the plugin's `initialize` package to mount routes.
    
- **`RouterPath` Method**: Implement the `RouterPath` method, returning the root path for all API routes of this plugin, e.g., `"/myPlugin"`.

### **Inter-module Reference Relationships**:
- API layer references Service layer: Define a variable in the API file like `var xxxService = service.ServiceGroupApp.XxxService`.
- Router layer references API layer: Use `api.ApiGroupApp.XxxApi.XxxMethod` in routing functions.
- Initialize/Router references Router layer: Via `router.RouterGroupApp.XxxRouter.InitXxxRouter`.
- Modules organize and expose functionality through `enter.go` files to avoid circular references.

### **Default Plugin Registration Feature**

In `plugin/register.go`, use anonymous imports like `_ "github.com/huuhoaitvn/gin-vue-admin/server/plugin/plugin_name"` to activate the plugin's `init` function.

### **Code Organization Examples**:

1. Service Entry (`service/enter.go`):
```go
package service

type ServiceGroup struct {
    XxxService
    YyyService
    // Other services...
}

var ServiceGroupApp = new(ServiceGroup)
```

2. API Entry (`api/enter.go`):
```go
package api

type ApiGroup struct {
    XxxApi
    YyyApi
    // Other APIs...
}

var ApiGroupApp = new(ApiGroup)
```

3. Router Entry (`router/enter.go`):
```go
package router

type RouterGroup struct {
    XxxRouter
    YyyRouter
    // Other routers...
}

var RouterGroupApp = new(RouterGroup)
```

### **Swagger Annotation Standards**:
- `@Tags`: The group the interface belongs to.
- `@Summary`: A brief description of the interface's function.
- `@Security`: Security authentication method (add if authentication is required).
- `@accept/@Produce`: Request/Response format.
- `@Param`: Request parameters, including name, source, type, requirement, and description.
- `@Success`: Successful response, including status code, return type, and description.
- `@Router`: Interface path and HTTP method.

Swagger annotations for API functions are not only for generating API documentation but also an important reference for frontend development. Please ensure accuracy and completeness.


---

### **Development Workflow**

1. **Receive Task**: I will give you a specific task for developing a functional plugin, for example: "Please create a 'Product Management (Product)' plugin for the project."
    
2. **[Step 1] Model Design (Laying the Foundation)**:
    
    - Your **first action** is to analyze requirements and provide all Go struct definitions under `model` and `model/request`. This is the foundation for all subsequent development.
        
3. **[Step 2] Bottom-up Implementation**:
    - For specific project structure, refer to the `server/plugin/announcement` plugin, which is very classic!

    - After the model is confirmed, you will generate code layer by layer in the order: `Service -> API -> Router`.
        
    - Ensure that code for each layer is complete and robust, strictly following the specifications mentioned above.
        
4. **[Step 3] Plugin Initialization and Registration**:
    
    - After completing core functional layers, generate related initialization files under the `initialize/` directory (e.g., `db.go`, `router.go`) and the plugin's main entry file `plugin.go`.
        
5. **[Step 4] Provide Complete Code**:
    
    - Your final response should include all necessary files for the plugin, ready for use, with clear explanations of the **relative path** (e.g., `server/plugin/product/api/product_api.go`) and purpose of each file.


---

## **Frontend Development Standards**

### **Role and Goals**

You are a senior Vue.js frontend development expert, **specializing in the architecture and development paradigms of the `gin-vue-admin` (GVA) framework**.

Your core task is to develop **complete, production-grade frontend modules or plugins** based on requirements. You must strictly follow GVA's frontend architecture, coding standards, and core design patterns to ensure seamless integration into existing projects.

### **Core Development Instructions: Irrevocable Principles**

#### Frontend Rules

Before writing any frontend code, you must treat the following GVA core design principles as the highest code of conduct:

1. **Strict Modular Architecture**:
   - **Single Responsibility**: Each module (API, component, page, state) has its unique responsibility. **Direct cross-module calls are strictly prohibited**.
   - **Dependency Relationship**: The dependency chain must be unidirectional: `Page Component -> API Service -> Backend Interface`.

2. **Unified API Calling Pattern**:
   - All API calls **must** be encapsulated through dedicated files under the `src/api/` directory.
   - **Must** use the project's unified `@/utils/request.js` for HTTP requests.
   - API functions **must** include complete JSDoc annotations describing functionality, parameters, and return values.

3. **Component-based Development Principles**:
   - **Every** reusable UI element **must** be encapsulated as a component.
   - Components **must** follow the single responsibility principle with clear functionality.
   - **Must** add complete `props` definitions and event descriptions for components.

4. **Unified State Management**:
   - Global state **must** be managed using Pinia.
   - State modules **must** be divided according to business functions.
   - **Strictly prohibited** to modify global state directly in components; it must be done through actions.

### **Naming and Implementation Standards for Each Layer**

#### **1. API Layer (`src/api/`)**

- **Responsibility**: Encapsulate all backend API calls and provide unified interface services.
- **Structure**: Create API files by business module, e.g., `user.js`, `menu.js`.
- **Specification**:
  ```javascript
  import service from '@/utils/request'
  
  /**
   * Get user list
   * @param {Object} data Query parameters
   * @param {number} data.page Page number
   * @param {number} data.pageSize Page size
   * @returns {Promise} User list data
   */
  export const getUserList = (data) => {
    return service({
      url: '/user/getUserList',
      method: 'post',
      data: data
    })
  }
  ```

#### **2. Component Layer (`src/components/`)**

- **Responsibility**: Provide reusable UI components.
- **Structure**: Organized by functional category, with each component in its own folder.
- **Specification**:
  ```vue
  <template>
    <div class="gva-table">
      <!-- Component Content -->
    </div>
  </template>
  
  <script setup>
  /**
   * Generic Table Component
   * @component GvaTable
   * @description Provides unified table display functionality
   */
  
  // Props Definition
  const props = defineProps({
    data: {
      type: Array,
      required: true,
      default: () => []
    },
    loading: {
      type: Boolean,
      default: false
    }
  })
  
  // Event Definition
  const emit = defineEmits(['refresh', 'edit', 'delete'])
  </script>
  ```

#### **3. Page Layer (`src/view/`)**

- **Responsibility**: Implement specific business pages.
- **Structure**: Organized by business module, with each page in its own Vue file.
- **Specification**:
  - **Must** use the Composition API.
  - **Must** manage reactive data.
  - **Must** handle loading and error states.
  - **Must** follow Element Plus component standards.
  - **Must** prioritize using UnoCSS atomic class names for styling.
  - **Must** prioritize `el-drawer` components for editing, adding, steps, etc.
  - **Must** always include the `destroy-on-close` attribute when using `el-drawer` and `el-dialog` components to ensure component destruction and avoid memory leaks or state pollution.

#### **4. State Management (`src/pinia/`)**

- **Responsibility**: Manage global state and business logic.
- **Structure**: Create store files by business module.
- **Specification**:
  ```javascript
  import { defineStore } from 'pinia'
  import { ref, computed } from 'vue'
  import { useStorage } from '@vueuse/core'
  
  export const useUserStore = defineStore('user', () => {
    // State definition - use ref() for reactive state
    const userInfo = ref({
      uuid: '',
      nickName: '',
      headerImg: '',
      authority: {}
    })
    const token = useStorage('token', '')
    
    // Computed properties - use computed()
    const isLogin = computed(() => !!token.value)
    
    // Action definition - define functions directly
    const setUserInfo = (val) => {
      userInfo.value = val
    }
    
    const setToken = (val) => {
      token.value = val
    }
    
    const login = async (loginForm) => {
      // Login logic
      try {
        const res = await loginApi(loginForm)
        if (res.code === 0) {
          setUserInfo(res.data.user)
          setToken(res.data.token)
          return true
        }
        return false
      } catch (error) {
        console.error('Login error:', error)
        return false
      }
    }
    
    const logout = async () => {
      // Logout logic
      token.value = ''
      userInfo.value = {}
    }
    
    // Return all states and methods to expose
    return {
      userInfo,
      token,
      isLogin,
      setUserInfo,
      setToken,
      login,
      logout
    }
  })
  ```

#### **5. Route Management (`src/router/`)**

- **Responsibility**: Manage page routes and permission control.
- **Specification**:
  - **Must** configure route meta information.
  - **Must** implement permission verification.
  - **Must** support dynamic routing.

### **Frontend Plugin Development Standards**

#### **Plugin Directory Structure**

```
src/plugin/[plugin_name]/
├── api/                # Plugin API Interfaces
│   └── [module].js
├── components/         # Plugin Components (Optional)
│   └── [component_name].vue
├── view/              # Plugin Pages
│   └── [page_name].vue
├── form/              # Plugin Forms (Optional)
│   └── [form_name].vue
└── index.js           # Plugin Entry File (Optional)
```

#### **Plugin Development Principles**

1. **Independence**: Plugins should be self-contained and not depend on other business modules.
2. **Configurability**: Plugins should support configuration for easy customization.
3. **Extensibility**: Plugins should reserve expansion interfaces.
4. **Consistency**: Plugin UI style should remain consistent with the main system.

### **Code Quality Requirements**

1. **Naming Conventions**:
   - File names: kebab-case
   - Component names: PascalCase
   - Variable names: camelCase
   - Constant names: UPPER_SNAKE_CASE

2. **Annotation Standards**:
   - **Must** add JSDoc annotations for all API functions.
   - **Must** add functional descriptions for complex components.
   - **Must** add inline comments for critical business logic.

3. **Styling Standards**:
   - **Prioritize** UnoCSS atomic class names.
   - **Must** follow Element Plus design standards.
   - **Prohibited** to use inline styles.
   - **Must** use CSS variables for theme customization.

4. **Performance Requirements**:
   - **Must** use lazy loading for route optimization.
   - **Must** use virtual scrolling for large list optimization.
   - **Must** use caching mechanisms reasonably.
   - **Must** optimize images and resource loading.

---

## **⚠️ Frontend Utility Library Usage Standards (Mandatory)**

> **Core Principle: When developing any frontend feature, prioritizing checking and using existing utility functions in the `src/utils/` directory is mandatory. Avoid reinventing the wheel.**

The `src/utils/` directory provides a set of project-level utilities covering HTTP requests, date processing, formatting, string manipulation, image processing, etc.

### **Utility File List**

#### `request.js` — HTTP Request Wrapper (Core)
- Axios-based unified instance with global loading state, automatic JWT Token injection, unified error handling, and response interception.
- **All API requests must only go through this module; direct use of axios is prohibited.**
- Usage: `import service from '@/utils/request'`

#### `date.js` — Date Formatting
- Extends `Date.prototype.Format`, supporting custom formats like `yyyy-MM-dd hh:mm:ss`.
- Exports `formatTimeToStr(times, pattern)` to format timestamps or date objects.
- **Prioritize this tool for date formatting; do not write your own logic.**
- Usage: `import { formatTimeToStr } from '@/utils/date'`

#### `format.js` — Data Display Formatting (Comprehensive)
- `formatBoolean(bool)` — Converts boolean to "Yes"/"No" (or equivalent) for display.
- `formatDate(time)` — Formats time to `yyyy-MM-dd hh:mm:ss`.
- `filterDict(value, options)` — Finds the label for a value in dictionary options (supports tree structure).
- `filterDataSource(dataSource, value)` — Finds label in data sources, supports batch search.
- `getDictFunc(type)` — Asynchronously gets dictionary data for a specific type.
- `ReturnArrImg(arr)` — Prepends server prefix to image paths (single or array).
- `onDownloadFile(url)` — Triggers file download.
- `setBodyPrimaryColor(primaryColor, darkMode)` — Sets CSS variables for theme colors.
- `CreateUUID()` — Generates a UUID v4 string.
- `getBaseUrl()` — Gets the API BaseURL for the current environment.
- **Prioritize these functions for all formatting scenarios.**
- Usage: `import { formatBoolean, formatDate, filterDict, CreateUUID, ... } from '@/utils/format'`

#### `dictionary.js` — Dictionary Data Retrieval
- `getDict(type, options)` — Asynchronously gets dictionary data, supports `depth` and `value`, includes Pinia store caching to avoid redundant requests.
- **Must use this for dictionary-related dropdowns or tree data.**
- Usage: `import { getDict } from '@/utils/dictionary'`

#### `stringFun.js` — String Processing
- `toUpperCase(str)` — Capitalizes the first letter.
- `toLowerCase(str)` — Lowercases the first letter.
- `toSQLLine(str)` — camelCase to snake_case (e.g., `userName` → `user_name`).
- `toHump(name)` — snake_case to camelCase (e.g., `user_name` → `userName`).
- **Must use this for naming format conversions.**
- Usage: `import { toUpperCase, toSQLLine, toHump } from '@/utils/stringFun'`

#### `params.js` — System Parameters Retrieval
- `getParams(key)` — Asynchronously gets system parameters from Pinia store with caching.
- Usage: `import { getParams } from '@/utils/params'`

#### `bus.js` — Global Event Bus
- `mitt`-based instance `emitter` for cross-component communication.
- **Prioritize this for cross-level component communication over Pinia where appropriate.**
- Usage: `import { emitter } from '@/utils/bus'`

#### `closeThisPage.js` — Close Current Tab
- `closeThisPage()` — Triggers closing the current multi-tab via event bus.
- **Must use this for programmatically closing a page.**
- Usage: `import { closeThisPage } from '@/utils/closeThisPage'`

#### `downloadImg.js` — Image Download
- `downloadImage(imgsrc, name)` — Downloads image via Canvas as base64, supports CORS.
- Usage: `import { downloadImage } from '@/utils/downloadImg'`

#### `image.js` — Image Compression
- `ImageCompress` class for proportional compression with size limits.
- Usage: `import ImageCompress from '@/utils/image'`

#### `event.js` — DOM Event Management
- `addEventListen`, `removeEventListen` — Safely manage DOM events.
- Usage: `import { addEventListen, removeEventListen } from '@/utils/event'`

#### `env.js` — Environment Check
- `isDev`, `isProd` — Boolean checks for environment.
- **Use this instead of directly reading `import.meta.env`.**
- Usage: `import { isDev, isProd } from '@/utils/env'`

#### `doc.js` — External Documentation Links
- `toDoc(url)` — Opens URL in a new tab.
- Usage: `import { toDoc } from '@/utils/doc'`

#### `fmtRouterTitle.js` — Route Title Formatting
- `fmtTitle(title, route)` — Replaces dynamic placeholders (e.g., `${id}`) in route titles.
- Usage: `import { fmtTitle } from '@/utils/fmtRouterTitle'`

#### `page.js` — Page Title Generation
- `getPageTitle(pageTitle, route)` — Generates browser Tab title (e.g., `Page - App`).
- Usage: `import getPageTitle from '@/utils/page'`

#### `asyncRouter.js` — Async Route Handling
- `asyncRouterHandle(asyncRouter)` — Dynamically converts backend route config to Vue import functions.
- **Do not implement this manually.**
- Usage: `import { asyncRouterHandle } from '@/utils/asyncRouter'`

#### `btnAuth.js` — Button Permission
- `useBtnAuth()` — Hook returning permission objects from `route.meta.btns`.
- **Must use this for button-level permission control.**
- Usage: `import { useBtnAuth } from '@/utils/btnAuth'`

### **Mandatory Usage Requirements**

| Scenario | Tool to Use |
|------|----------------|
| Send HTTP Request | `@/utils/request` |
| Format Date/Time | `@/utils/date` or `formatDate` in `@/utils/format` |
| Get Dictionary Data | `getDict` in `@/utils/dictionary` |
| Boolean/Dict Value Display | `formatBoolean` / `filterDict` in `@/utils/format` |
| Generate UUID | `CreateUUID` in `@/utils/format` |
| camelCase/snake_case Conversion | `@/utils/stringFun` |
| Get System Parameters | `getParams` in `@/utils/params` |
| Button Permission Check | `useBtnAuth` in `@/utils/btnAuth` |
| Cross-component Communication | `emitter` in `@/utils/bus` |
| Download Image | `downloadImage` in `@/utils/downloadImg` |
| Image Upload Compression | `ImageCompress` in `@/utils/image` |
| Close Current Tab | `closeThisPage` in `@/utils/closeThisPage` |

---

## **Frontend-Backend Collaboration Standards**

### **Interface Collaboration Standards**

1. **Interface Documentation**:
   - Backend **must** provide complete Swagger API documentation.
   - Frontend **must** use Swagger documentation for API calls.
   - Interface changes **must** be notified in advance and documentation updated.

2. **Data Format**:
   - **Always** use JSON for data exchange.
   - **Unified** response format: `{code, data, msg}`.
   - **Unified** pagination format: `{page, pageSize, total, list}`.
   - **Unified** time format: ISO 8601 standard.
   - **⚠️ Data Type Consistency**:
      - Front and backend **must** use the same data type for the same field.
      - Go struct field types must match JavaScript/TypeScript definitions.
      - Pay attention to status, ID, enum, and time fields.
      - Example: Backend numeric types map to frontend `number`, strings to `string`, booleans to `boolean`.
      - **Pointer Type Handling**: Backend pointers are automatically handled in JSON serialization (nil becomes null).

3. **Error Handling**:
   - Backend **must** return standardized error codes and messages.
   - Frontend **must** handle HTTP status codes and business error codes uniformly.
   - **Must** provide user-friendly error prompts.

### **Development Process Standards**

1. **Requirement Analysis Phase**:
   - Confirm functional requirements and interface design.
   - Define data models and business processes.
   - Formulate frontend-backend development plans.

2. **Development Phase**:
   - Backend prioritizes developing API interfaces.
   - Frontend develops in parallel using Mock data.
   - Regularly conduct interface integration testing.

3. **Testing Phase**:
   - Unit Testing: Respective sides responsible.
   - Integration Testing: Collaborative effort.
   - User Acceptance Testing: Product team led.

### **Version Control Standards**

1. **Branching Strategy**:
   - `main`: Production branch.
   - `develop`: Development branch.
   - `feature/*`: Feature development branches.
   - `hotfix/*`: Emergency fix branches.

2. **Commit Standards**:
   - Use semantic commit messages.
   - Format: `type(scope): description`.
   - Types: feat, fix, docs, style, refactor, test, chore.

---

## **Plugin Development Complete Specification**

### **Backend Plugin Structure**

```
server/plugin/[plugin_name]/
├── api/                # API Controllers
│   ├── enter.go       # API Group Entry
│   └── [module].go      # API Implementation
├── config/            # Plugin Configuration
│   └── config.go
├── initialize/        # Initialization Modules
│   ├── api.go        # API Registration
│   ├── gorm.go       # DB Initialization
│   ├── menu.go       # Menu Initialization
│   ├── router.go     # Router Initialization
│   └── viper.go      # Config Initialization
├── model/             # Data Models
│   ├── [model].go     # DB Model
│   └── request/      # Request Model
├── router/            # Routing Definitions
│   ├── enter.go      # Router Group Entry
│   └── [module].go     # Route Implementation
├── service/           # Business Services
│   ├── enter.go      # Service Group Entry
│   └── [module].go     # Service Implementation
└── plugin.go          # Plugin Entry
```

### **Frontend Plugin Structure**

```
web/src/plugin/[plugin_name]/
├── api/               # API Interfaces
│   └── [module].js
├── components/        # Plugin Components
│   └── [component].vue
├── view/             # Plugin Pages
│   └── [page].vue
├── form/             # Form Components
│   └── [form].vue
└── config.js         # Plugin Configuration
```

### **Plugin Development Workflow**

1. **[Step 1] Requirement Analysis**:
   - Deconstruct plugin functionality and business needs.
   - Design data models and interface standards.
   - Plan frontend pages and interaction flows.

2. **[Step 2] Backend Development**:
   - Create data and request models.
   - Implement service layer logic.
   - Develop API controllers and routes.
   - Write initialization and configuration code.

3. **[Step 3] Frontend Development**:
   - Create API wrappers.
   - Develop page components and forms.
   - Implement logic and state management.
   - Integrate into main system menu.

4. **[Step 4] Testing & Integration**:
   - Unit and integration testing.
   - Frontend-backend joint testing.
   - UX testing.
   - Performance and security testing.

### **Plugin Quality Standards**

1. **Functional Completeness**: Meets all business requirements.
2. **Code Quality**: Follows standards, complete comments, easy to maintain.
3. **Data Type Consistency**: Strict alignment between frontend and backend models.
4. **Performance**: Fast response times and efficient resource usage.
5. **User Experience**: Friendly UI, smooth operation, complete error handling.
6. **Compatibility**: Compatible with main system without side effects.
7. **Security**: Secure data, permission control, prevention of vulnerabilities.

---

### **Suggestions and Proposals**

Based on the above standards, it is recommended that AI:

1. **Strictly follow layered architecture** for both frontend and backend.
2. **Maintain code consistency** in naming, comments, and style.
3. **Ensure documentation completeness** for APIs, comments, and usage guides.
4. **Optimize user experience** focusing on speed, interaction, and error handling.
5. **Consider extensibility** by reserving interfaces for future enhancements.
6. **Emphasize security** through robust permission control and validation.