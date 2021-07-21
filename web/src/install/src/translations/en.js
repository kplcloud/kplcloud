export default {
  common: {
    add: 'Add',
    cancel: 'Cancel',
    description: 'Description',
    delete: 'Delete',
    title: 'Title',
    save: 'Save',
    faq: 'FAQ',
    contact: 'Contact Us',
    tos: 'Terms of Service',
    policy: 'Privacy Policy'
  },
  board: {
    titlePlaceholder: 'Enter a title for this card',
    deleteDescription: 'Are you sure you want to delete this card?',
    editCard: 'Edit Card',
    deleteCard: 'Delete Card',
    state: {
      TODO: 'TO DO',
      INPROGRESS: 'INPROGRESS',
      TESTING: 'TESTING',
      DONE: 'DONE'
    }
  },
  chat: {
    online: 'Users Online ({count})',
    addChannel: 'Add Channel',
    channel: 'Channel | Channels',
    message: 'Message'
  },
  email: {
    compose: 'Compose Email',
    send: 'Send',
    subject: 'Subject',
    labels: 'Labels',
    emptyList: 'Empty email list',
    inbox: 'Inbox',
    sent: 'Sent',
    drafts: 'Drafts',
    starred: 'Starred',
    trash: 'Trash',
    work: 'Work',
    invoice: 'Invoice'
  },
  todo: {
    addTask: 'Add Task',
    tasks: 'Tasks',
    completed: 'Completed',
    labels: 'Labels'
  },
  dashboard: {
    activity: 'Activity',
    weeklySales: 'Weekly Sales',
    sales: 'Sales',
    recentOrders: 'Recent Orders',
    sources: 'Traffic Sources',
    lastweek: 'vs last week',
    orders: 'Orders',
    customers: 'Customers',
    tickets: 'Support Tickets',
    viewReport: 'View Report'
  },
  usermenu: {
    profile: 'Profile',
    signin: 'Sign In',
    dashboard: 'Dashboard',
    signout: 'Sign Out'
  },
  error: {
    notfound: 'Page Not Found',
    other: 'An Error Ocurred'
  },
  check: {
    title: 'Set New Password',
    backtosign: 'Back to Sign In',
    newpassword: 'New Password',
    button: 'Set new password and Sign in',
    error: 'The action link is invalid',
    verifylink: 'Verifying link...',
    verifyemail: 'Verifying email address...',
    emailverified: 'Email verified! Redirecting...'
  },
  forgot: {
    title: 'Forgot Password?',
    subtitle: 'Enter your account email address and we will send you a link to reset your password.',
    email: 'Email',
    button: 'Request Password Reset',
    backtosign: 'Back to Sign In'
  },
  login: {
    title: 'Sign In',
    email: 'Email',
    password: 'Password',
    button: 'Sign In',
    orsign: 'Or sign in with',
    forgot: 'Forgot password?',
    noaccount: 'Don\'t have an account?',
    create: 'Create one here',
    error: 'The email / password combination is invalid',
    loginaccount: 'Log in to your account'
  },
  register: {
    title: 'Create Account',
    name: 'Full name',
    email: 'Email',
    password: 'Password',
    button: 'Create Account',
    orsign: 'Or sign up with',
    agree: 'By signing up, you agree to the',
    account: 'Already have an account?',
    signin: 'Sign In'
  },
  utility: {
    maintenance: 'In Maintenance'
  },
  faq: {
    call: 'Have other questions? Please reach out '
  },
  ecommerce: {
    products: 'Products',
    filters: 'Filters',
    collections: 'Collections',
    priceRange: 'Price Range',
    customerReviews: 'Customer Reviews',
    up: 'and up',
    brand: 'Brand',
    search: 'Search for product',
    results: 'Results ( {0} of {1} )',
    orders: 'Orders',
    shipping: 'Shipping',
    freeShipping: 'Free Shipping',
    inStock: 'In Stock',
    quantity: 'Quantity',
    addToCart: 'Add To Cart',
    buyNow: 'Buy Now',
    price: 'Price',
    about: 'About this item',
    description: 'Description',
    reviews: 'Reviews',
    details: 'Product Details',
    cart: 'Cart',
    summary: 'Order Summary',
    total: 'Total',
    discount: 'Discount',
    subtotal: 'Subtotal',
    continue: 'Continue Shopping',
    checkout: 'Checkout'
  },
  menu: {
    search: 'Search (press "ctrl + /" to focus)',
    dashboard: 'Dashboard',
    logout: 'Logout',
    profile: 'Profile',
    blank: 'Blank Page',
    pages: 'Pages',
    others: 'Others',
    email: 'Email',
    chat: 'Chat',
    todo: 'To Do',
    board: 'Task Board',
    users: 'Users',
    usersList: 'List',
    usersEdit: 'Edit',
    ecommerce: 'Ecommerce',
    ecommerceList: 'Products',
    ecommerceProductDetails: 'Product Details',
    ecommerceOrders: 'Orders',
    ecommerceCart: 'Cart',
    auth: 'Auth Pages',
    authLogin: 'Signin / Login',
    authRegister: 'Signup / Register',
    authVerify: 'Verify Email',
    authForgot: 'Forgot Password',
    authReset: 'Reset Password',
    errorPages: 'Error Pages',
    errorNotFound: 'Not Found / 404',
    errorUnexpected: 'Unexpected / 500',
    utilityPages: 'Utility Pages',
    utilityMaintenance: 'Maintenance',
    utilitySoon: 'Coming Soon',
    utilityHelp: 'FAQs / Help',
    levels: 'Menu Levels',
    disabled: 'Menu Disabled',
    docs: 'Documentation',
    feedback: 'Feedback',
    support: 'Support',
    landingPage: 'Landing Page',
    pricingPage: 'Pricing Page',
    systemPage: 'System Manage',
    systemSettingPage: 'System Setting'
  },
  // Vuetify components translations
  $vuetify: {
    badge: 'Badge',
    close: 'Close',
    dataIterator: {
      noResultsText: 'No matching records found',
      loadingText: 'Loading items...'
    },
    dataTable: {
      itemsPerPageText: 'Rows per page:',
      ariaLabel: {
        sortDescending: 'Sorted descending.',
        sortAscending: 'Sorted ascending.',
        sortNone: 'Not sorted.',
        activateNone: 'Activate to remove sorting.',
        activateDescending: 'Activate to sort descending.',
        activateAscending: 'Activate to sort ascending.'
      },
      sortBy: 'Sort by'
    },
    dataFooter: {
      itemsPerPageText: 'Items per page:',
      itemsPerPageAll: 'All',
      nextPage: 'Next page',
      prevPage: 'Previous page',
      firstPage: 'First page',
      lastPage: 'Last page',
      pageText: '{0}-{1} of {2}'
    },
    datePicker: {
      itemsSelected: '{0} selected',
      nextMonthAriaLabel: 'Next month',
      nextYearAriaLabel: 'Next year',
      prevMonthAriaLabel: 'Previous month',
      prevYearAriaLabel: 'Previous year'
    },
    noDataText: 'No data available',
    carousel: {
      prev: 'Previous visual',
      next: 'Next visual',
      ariaLabel: {
        delimiter: 'Carousel slide {0} of {1}'
      }
    },
    calendar: {
      moreEvents: '{0} more'
    },
    fileInput: {
      counter: '{0} files',
      counterSize: '{0} files ({1} in total)'
    },
    timePicker: {
      am: 'AM',
      pm: 'PM'
    },
    pagination: {
      ariaLabel: {
        wrapper: 'Pagination Navigation',
        next: 'Next page',
        previous: 'Previous page',
        page: 'Goto Page {0}',
        currentPage: 'Current Page, Page {0}'
      }
    }
  }
}
