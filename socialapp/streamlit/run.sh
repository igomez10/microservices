#!/bin/bash

# Socialapp Admin - Quick Start Script

echo "🔧 Socialapp Admin Tool"
echo "======================="
echo ""

# Check if we're in the right directory
if [ ! -f "app.py" ]; then
    echo "❌ Error: app.py not found"
    echo "Please run this script from the streamlit directory"
    exit 1
fi

# Check if Python is installed
if ! command -v python3 &> /dev/null; then
    echo "❌ Error: Python 3 is not installed"
    exit 1
fi

# Check if streamlit is installed
if ! python3 -c "import streamlit" &> /dev/null; then
    echo "⚠️  Streamlit not found. Installing dependencies..."
    pip3 install -r requirements.txt
fi

echo "✅ Starting Socialapp Admin..."
echo ""

# Run streamlit
streamlit run app.py

