<template>
  <div class="min-h-screen bg-gray-50">
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex">
            <div class="flex-shrink-0 flex items-center">
              <h1 class="text-xl font-bold text-gray-800">TicketHub</h1>
            </div>
            <nav class="ml-6 flex space-x-4">
              <NuxtLink to="/" class="px-3 py-2 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-100">
                Dashboard
              </NuxtLink>
              <NuxtLink to="/issues" class="px-3 py-2 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-100">
                Issues
              </NuxtLink>
            </nav>
          </div>
        </div>
      </div>
    </header>

    <main>
      <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <NuxtPage />
      </div>
    </main>

    <footer class="bg-white">
      <div class="max-w-7xl mx-auto px-4 py-6 sm:px-6 lg:px-8">
        <p class="text-center text-sm text-gray-500">
          &copy; 2025 TicketHub -
        </p>
      </div>
    </footer>
  </div>
</template>
