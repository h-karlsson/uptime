import { FlaskConical } from 'lucide-react'
import { Badge } from '@/components/ui/badge'

export function MockModeBadge() {
  return (
    <Badge variant="secondary" className="fixed bottom-4 left-4 z-50 gap-1.5 shadow-sm">
      <FlaskConical />
      Mock data
    </Badge>
  )
}
