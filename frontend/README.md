# CityStars Frontend

A mobile-first, responsive web application for the CityStars civic reporting platform.

## 📁 Project Structure

```
frontend/
├── index.html                  # Homepage
├── pages/                      # All application pages
│   ├── login.html             # User login
│   ├── register.html          # User registration
│   ├── create-report.html     # Create new report
│   ├── reports.html           # Browse all reports
│   ├── report-detail.html     # View single report
│   ├── my-reports.html        # User's reports
│   ├── leaderboard.html       # Top contributors
│   ├── stats.html             # Statistics dashboard
│   └── profile.html           # User profile
├── assets/
│   ├── css/                   # Stylesheets
│   │   ├── main.css          # Core styles & variables
│   │   ├── components.css    # Component styles
│   │   ├── forms.css         # Form styles
│   │   └── responsive.css    # Media queries
│   └── js/                    # JavaScript files
│       ├── config.js         # API configuration
│       ├── utils.js          # Utility functions
│       ├── api.js            # API service layer
│       ├── auth.js           # Authentication management
│       ├── main.js           # Homepage logic
│       └── pages/            # Page-specific scripts
│           ├── login.js
│           ├── register.js
│           ├── create-report.js
│           ├── reports.js
│           ├── report-detail.js
│           ├── my-reports.js
│           ├── leaderboard.js
│           ├── stats.js
│           └── profile.js
```

## 🚀 Features

### User Features
- **Authentication**: Register and login with phone number
- **Create Reports**: Submit civic issues with images
- **Browse Reports**: Filter by status and category
- **Track Reports**: View your submitted reports
- **Leaderboard**: See top contributors
- **Statistics**: View community impact

### Admin Features
- **Update Reports**: Change status (pending → in-progress → completed/rejected)
- **Add After Images**: Upload completion proof
- **Full Access**: View and manage all reports

## 🎨 Design Features

- **Mobile-First**: Optimized for mobile devices
- **Responsive**: Adapts to all screen sizes
- **Modern UI**: Clean, intuitive interface
- **Toast Notifications**: User feedback
- **Loading States**: Visual feedback for async operations
- **Error Handling**: Comprehensive validation and error messages

## 🔧 Configuration

Edit `assets/js/config.js` to configure the API endpoint:

```javascript
const API_CONFIG = {
    BASE_URL: 'http://localhost:4000',  // Change to your API URL
    VERSION: 'v1'
};
```

## 📱 Pages Overview

### Public Pages
- **Homepage** (`index.html`): Landing page with stats and recent reports
- **Reports** (`pages/reports.html`): Browse all reports with filters
- **Leaderboard** (`pages/leaderboard.html`): Top 10 contributors
- **Stats** (`pages/stats.html`): Community statistics

### Authentication Pages
- **Login** (`pages/login.html`): User login
- **Register** (`pages/register.html`): User registration

### Protected Pages (Require Login)
- **Create Report** (`pages/create-report.html`): Submit new report
- **My Reports** (`pages/my-reports.html`): User's submitted reports
- **Profile** (`pages/profile.html`): User profile and stats
- **Report Detail** (`pages/report-detail.html`): View single report with admin actions

## 🎯 API Integration

The frontend integrates with your Go backend API:

- **POST** `/v1/user/register` - User registration
- **POST** `/v1/user/login` - User authentication
- **GET** `/v1/user/me` - Get user profile
- **POST** `/v1/reports` - Create report
- **GET** `/v1/reports` - List all reports
- **GET** `/v1/reports/{id}` - Get report details
- **PATCH** `/v1/reports/{id}` - Update report (admin)
- **GET** `/v1/user/reports` - Get user's reports
- **GET** `/v1/reports/stats` - Get statistics
- **GET** `/v1/leaderboard` - Get leaderboard

## 🔐 Authentication

- JWT tokens stored in localStorage
- Automatic token validation
- Session expiry handling
- Role-based access control (user/admin)

## 🌐 Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)
- Mobile browsers (iOS Safari, Chrome Mobile)

## 📝 Usage

1. **Start your Go backend** (default: `http://localhost:4000`)
2. **Open `index.html`** in a web browser or serve with a local server
3. **Register** a new account or **login**
4. **Start reporting** civic issues!

## 🎨 Customization

### Colors
Edit CSS variables in `assets/css/main.css`:
```css
:root {
    --primary-color: #6366f1;
    --secondary-color: #10b981;
    /* ... more variables */
}
```

### Categories
Edit category icons in `assets/js/config.js`:
```javascript
const CATEGORY_ICONS = {
    pothole: '🕳️',
    streetlight: '💡',
    // ... add more
};
```

## 🔨 Development

For development with live reload, use a local server:

```bash
# Python
python -m http.server 8000

# Node.js
npx http-server

# PHP
php -S localhost:8000
```

Then visit `http://localhost:8000`

## 🐛 Known Issues

- Image upload uses URLs (not file upload) - update if you implement file upload
- No image optimization - large images may load slowly
- Session management relies on localStorage - consider more secure alternatives for production

## 📄 License

Same as the CityStars backend project.
