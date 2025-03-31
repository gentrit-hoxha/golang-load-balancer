#!/usr/bin/env python3

import sys
import json
from collections import defaultdict

def analyze_logs(log_file):
    results = {
        "roles": {
            "admin": {"total_requests": 0, "success_count": 0, "success_rate": 0},
            "client": {"total_requests": 0, "success_count": 0, "success_rate": 0},
            "user": {"total_requests": 0, "success_count": 0, "success_rate": 0}
        },
        "summary": {
            "total_requests": 0,
            "successful_requests": 0,
            "failed_requests": 0,
            "overall_success_rate": 0
        },
        "status_codes": defaultdict(int)
    }
    
    with open(log_file, 'r') as f:
        for line in f:
            line = line.strip()
            if not line or ":" not in line:
                continue
                
            try:
                # Parse role and status code
                role, status_code = line.split(":", 1)
                status_code = int(status_code)
                
                # Update role statistics
                if role in results["roles"]:
                    results["roles"][role]["total_requests"] += 1
                    
                    # Check if request was successful
                    if 200 <= status_code < 300:
                        results["roles"][role]["success_count"] += 1
                
                # Update status code counts
                results["status_codes"][str(status_code)] += 1
                
            except (ValueError, IndexError):
                # Skip invalid lines
                continue
    
    # Calculate totals and percentages
    total_requests = 0
    successful_requests = 0
    
    for role, data in results["roles"].items():
        total_requests += data["total_requests"]
        successful_requests += data["success_count"]
        
        # Calculate success rate for this role
        if data["total_requests"] > 0:
            data["success_rate"] = (data["success_count"] / data["total_requests"]) * 100
        else:
            data["success_rate"] = 0
    
    # Update summary statistics
    results["summary"]["total_requests"] = total_requests
    results["summary"]["successful_requests"] = successful_requests
    results["summary"]["failed_requests"] = total_requests - successful_requests
    
    if total_requests > 0:
        results["summary"]["overall_success_rate"] = (successful_requests / total_requests) * 100
    
    return results

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 analyze_logs.py <log_file>")
        print("Using default log file: load_test.log")
        log_file = "load_test.log"
    else:
        log_file = sys.argv[1]
    
    try:
        results = analyze_logs(log_file)
        
        # Print results in a readable format
        print("\n===== LOAD TEST RESULTS =====")
        print(f"Total Requests: {results['summary']['total_requests']}")
        print(f"Successful Requests: {results['summary']['successful_requests']}")
        print(f"Failed Requests: {results['summary']['failed_requests']}")
        print(f"Overall Success Rate: {results['summary']['overall_success_rate']:.2f}%")
        
        print("\nResults by Role:")
        for role, data in results["roles"].items():
            print(f"  {role.upper()}: {data['total_requests']} requests, {data['success_rate']:.2f}% success rate")
        
        print("\nStatus Code Distribution:")
        for status_code, count in sorted(results["status_codes"].items()):
            print(f"  {status_code}: {count} requests")
        
        # Also output as JSON
        print("\nJSON Output:")
        print(json.dumps(results, indent=2))
        
    except Exception as e:
        print(f"Error analyzing logs: {e}")
        sys.exit(1)

if __name__ == '__main__':
    main() 