import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'             the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/register',
    component: Layout,
    hidden: true,
    children: [{
      path: 'register',
      name: 'register',
      component: () => import('@/views/register/index'),
      meta: { title: 'I2P仪表盘', icon: 'dashboard' }
    }]
  },
  {
    path: '/detail',
    component: Layout,
    // hidden: true,
    children: [{
      path: 'detail',
      name: 'detail',
      component: () => import('@/views/form/experdetail'),
      meta: { title: 'detail', icon: 'dashboard' }
    }]
  },

  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },

  {
    path: '/403',
    component: () => import('@/views/403'),
    hidden: true
  },
  // {
  //   path: '/metadata',
  //   component: Layout,
  //   children: [
  //     {
  //       path: 'metadata',
  //       name: 'Metadata',
  //       component: () => import('@/views/metadata/index'),
  //       meta: { title: 'I2P仪表盘', icon: 'dashboard' }
  //     }
  //   ],
  //   // hidden: true
  // },

  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [{
      path: 'dashboard',
      name: 'Dashboard',
      component: () => import('@/views/dashboard/index'),
      meta: { title: 'I2P仪表盘', icon: 'dashboard' }
    }]
  },

  {
    path: '/example',
    component: Layout,
    redirect: '/example/table',
    name: 'Example',
    meta: { title: '元数据', icon: 'example' },
    children: [
      {
        path: 'table',
        name: 'Table',
        component: () => import('@/views/table/index'),
        meta: { title: '节点元数据', icon: 'table' }
      },
      {
        path: 'tree',
        name: 'Tree',
        component: () => import('@/views/tree/index'),
        meta: { title: 'Tree', icon: 'tree' },
        hidden: true
      }
    ]
  },
  {
    path: '/view',
    component: Layout,
    redirect: '/view/table',
    name: 'View',
    meta: { title: '定制可视化', icon: 'example' },
    children: [
      {
        path: 'table',
        name: 'Table',
        component: () => import('@/views/viewers/index'),
        meta: { title: '可视化定制', icon: 'example' }
      },
    ]
  },
  {
    path: '/pcap',
    component: Layout,
    redirect: '/pcap/analyser',
    name: 'View',
    meta: { title: '流量数据分析', icon: 'example' },
    children: [
      {
        path: 'pcap',
        name: 'pcap',
        component: () => import('@/views/pcapAnalyser/index'),
        meta: { title: 'Pcap-Analyser', icon: 'pdf' }
      },
      {
        path: 'pcap_metadata',
        name: 'metadata',
        component: () => import('@/views/pcapAnalyser/index_meta'),
        meta: { title: 'Pcap-Metadata', icon: 'pdf' }
      },
      {
        path: 'pcap_proto',
        name: 'proto',
        component: () => import('@/views/pcapAnalyser/index_proto'),
        meta: { title: 'Pcap-Proto', icon: 'pdf' }
      },
      {
        path: 'pcap_flow',
        name: 'flow',
        component: () => import('@/views/pcapAnalyser/index_flow'),
        meta: { title: 'Pcap-Flow', icon: 'pdf' }
      },
      {
        path: 'pcap_map',
        name: 'map',
        component: () => import('@/views/pcapAnalyser/index_pcap_map'),
        meta: { title: 'Pcap-Map', icon: 'pdf' }
      },
    ]
  },

  {
    path: '/form',
    component: Layout,
    redirect: '/form/index',
    meta: { title: '实验分析', icon: 'form' },
    name: 'LabNested',
    children: [
      {
        path: 'index',
        name: 'Form',
        component: () => import('@/views/form/experlist'),
        meta: { title: '实验列表', icon: 'form' },
      },
      {
        path: 'newExper',
        name: 'exper',
        component: () => import('@/views/form/index'),
        meta: { title: '新增实验', icon: 'form' },
      },

    ]
  },

  {
    path: '/nested',
    component: Layout,
    redirect: '/nested/menu1',
    name: 'Nested',
    meta: {
      title: 'probe',
      icon: 'nested'
    },
    children: [
      {
        path: 'menu1',
        component: () => import('@/views/probe/index'), // Parent router-view
        name: 'Menu1',
        meta: { title: '探针管理' },
      },
    ]
  },
  {
    path: '/nested',
    component: Layout,
    redirect: '/nested/menu1',
    name: 'Nested',
    meta: {
      title: 'Exper',
      icon: 'nested'
    },
    children: [
      {
        path: 'menu12',
        component: () => import('@/views/Exper/index2'), // Parent router-view
        name: 'Menu1',
        meta: { title: 'Exper' },
      },
    ]
  },
  // {
  //   path: '/nested',
  //   component: Layout,
  //   redirect: '/nested/menu1',
  //   name: 'Nested',
  //   meta: {
  //     title: 'probe',
  //     icon: 'nested'
  //   },
  //   children: [
  //     {
  //       path: 'menu1',
  //       component: () => import('@/views/nested/menu1/index'), // Parent router-view
  //       name: 'Menu1',
  //       meta: { title: 'Menu1' },
  //       children: [
  //         {
  //           path: 'menu1-1',
  //           component: () => import('@/views/nested/menu1/menu1-1'),
  //           name: 'Menu1-1',
  //           meta: { title: 'Menu1-1' }
  //         },
  //         {
  //           path: 'menu1-2',
  //           component: () => import('@/views/nested/menu1/menu1-2'),
  //           name: 'Menu1-2',
  //           meta: { title: 'Menu1-2' },
  //           children: [
  //             {
  //               path: 'menu1-2-1',
  //               component: () => import('@/views/nested/menu1/menu1-2/menu1-2-1'),
  //               name: 'Menu1-2-1',
  //               meta: { title: 'Menu1-2-1' }
  //             },
  //             {
  //               path: 'menu1-2-2',
  //               component: () => import('@/views/nested/menu1/menu1-2/menu1-2-2'),
  //               name: 'Menu1-2-2',
  //               meta: { title: 'Menu1-2-2' }
  //             }
  //           ]
  //         },
  //         {
  //           path: 'menu1-3',
  //           component: () => import('@/views/nested/menu1/menu1-3'),
  //           name: 'Menu1-3',
  //           meta: { title: 'Menu1-3' }
  //         }
  //       ]
  //     },
  //     {
  //       path: 'menu2',
  //       component: () => import('@/views/nested/menu2/index'),
  //       meta: { title: 'menu2' }
  //     }
  //   ]
  // },

  {
    path: 'external-link',
    component: Layout,
    children: [
      {
        path: 'https://panjiachen.github.io/vue-element-admin-site/#/',
        meta: { title: 'External Link', icon: 'link' }
      }
    ]
  },

  // 404 page must be placed at the end !!!
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
