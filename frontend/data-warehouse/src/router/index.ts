import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue'),
    },
    {
      path: '/proteins',
      name: 'protein-list',
      component: () => import('@/views/ProteinListView.vue'),
    },
    {
      path: '/proteins/new',
      name: 'protein-create',
      component: () => import('@/views/ProteinCreateView.vue'),
    },
    {
      path: '/proteins/:accession',
      name: 'protein-detail',
      component: () => import('@/views/ProteinDetailView.vue'),
      props: true,
    },
    {
      path: '/proteins/:accession/edit',
      name: 'protein-edit',
      component: () => import('@/views/ProteinEditView.vue'),
      props: true,
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/NotFoundView.vue'),
    },
  ],
})

export default router
